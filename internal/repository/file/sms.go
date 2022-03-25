package file

import (
	"context"
	"io/ioutil"
	"log"
	"sort"
	"strings"

	"antares-me/monitoring-system/internal/domain"

	"antares-me/monitoring-system/pkg/cache"
	"antares-me/monitoring-system/pkg/country"
	"antares-me/monitoring-system/pkg/provider"
)

type SmsRepo struct {
	file  string
	data  []domain.SMSData
	cache *cache.Cache
}

func NewSmsRepo(fp string, c *cache.Cache) *SmsRepo {
	return &SmsRepo{
		file:  fp,
		data:  []domain.SMSData{},
		cache: c,
	}
}

// parseData парсит данные из файла СМС в структуру репозитория
func (r *SmsRepo) parseData(ctx context.Context) error {
	r.data = []domain.SMSData{}
	content, err := ioutil.ReadFile(r.file)
	if err != nil {
		return err
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		delimitersCount := strings.Count(line, ";")
		if delimitersCount != 3 {
			log.Println("SMS: неверное количество полей:", delimitersCount+1)
			continue
		}
		fields := strings.Split(line, ";")
		if validateSmsFields(fields) {
			r.data = append(r.data, domain.SMSData{Country: fields[0], Bandwidth: fields[1], ResponseTime: fields[2], Provider: fields[3]})
		}
	}
	return nil
}

// validateSmsFields проверяет поля одной СМС записи на соответствие требованиям ТЗ.
func validateSmsFields(fields []string) bool {
	if !country.IsRightCode(fields[0]) {
		log.Println("SMS: неверный код страны:", fields[0])
		return false
	}
	if !provider.IsRightCode(fields[3], "Sms") {
		log.Println("SMS: неверный провайдер:", fields[3])
		return false
	}
	return true
}

// replaceCountryCodes заменяет коды страны на названия.
func (r *SmsRepo) replaceCountryCodes() {
	for i, sms := range r.data {
		r.data[i].Country = country.ReplaceCountruCode(sms.Country)
	}
}

// SortByProvider сортирует срез данных SMS по провайдеру.
func (r *SmsRepo) sortByProvider() []domain.SMSData {
	data := make([]domain.SMSData, len(r.data))
	copy(data, r.data)
	sort.SliceStable(data, func(i, j int) bool {
		return data[i].Provider < data[j].Provider
	})
	return data
}

// SortByProvider сортирует срез данных SMS по стране.
func (r *SmsRepo) sortByCountry() []domain.SMSData {
	data := make([]domain.SMSData, len(r.data))
	copy(data, r.data)
	sort.SliceStable(data, func(i, j int) bool {
		return data[i].Country < data[j].Country
	})
	return data
}

func (r *SmsRepo) GetResultData(ctx context.Context) ([][]domain.SMSData, error) {
	if val, has := r.cache.Get("sms"); has == true {
		v := val.([][]domain.SMSData)
		return v, nil
	} else {
		data := [][]domain.SMSData{}
		if err := r.parseData(ctx); err != nil {
			return data, err
		}
		r.replaceCountryCodes()
		data = append(data, r.sortByProvider())
		data = append(data, r.sortByCountry())
		r.cache.Set("sms", data, 0)
		return data, nil
	}
}
