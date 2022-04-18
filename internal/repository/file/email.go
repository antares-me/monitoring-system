package file

import (
	"context"
	"io/ioutil"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"

	"antares-me/monitoring-system/internal/domain"

	"antares-me/monitoring-system/pkg/cache"
	"antares-me/monitoring-system/pkg/country"
	"antares-me/monitoring-system/pkg/provider"
)

type EmailRepo struct {
	file  string
	data  []domain.EmailData
	cache *cache.Cache
}

func NewEmailRepo(fp string, c *cache.Cache) *EmailRepo {
	return &EmailRepo{
		file:  fp,
		data:  []domain.EmailData{},
		cache: c,
	}
}

func (r *EmailRepo) parseData(ctx context.Context) error {
	r.data = []domain.EmailData{}
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
		if delimitersCount != 2 {
			log.Println("Email: неверное количество полей:", delimitersCount+1)
			continue
		}
		fields := strings.Split(line, ";")
		if validateEmailFields(fields) {
			deliveryTime, err := strconv.Atoi(fields[2])
			if err != nil {
				continue
			}
			r.data = append(r.data, domain.EmailData{
				Country:      fields[0],
				Provider:     fields[1],
				DeliveryTime: deliveryTime,
			})
		}
	}
	return nil
}

// validateEmailFields проверяет поля одной Email записи на соответствие требованиям ТЗ
func validateEmailFields(fields []string) bool {
	if !country.IsRightCode(fields[0]) {
		log.Println("Email: неверный код страны:", fields[0])
		return false
	}
	if !provider.IsRightCode(fields[1], "Email") {
		log.Println("Email: неверный провайдер:", fields[1])
		return false
	}
	return true
}

func (r *EmailRepo) GetResultData(ctx context.Context, wg *sync.WaitGroup, res *domain.ResultSetT, e *[]error) {
	defer wg.Done()
	if val, has := r.cache.Get("email"); has == true {
		res.Email = val.(map[string][][]domain.EmailData)
	} else {
		data := make(map[string][][]domain.EmailData)
		if err := r.parseData(ctx); err != nil {
			*e = append(*e, err)
			return
		}
		byCountryMap := make(map[string][]domain.EmailData)
		for _, v := range r.data {
			byCountryMap[v.Country] = append(byCountryMap[v.Country], v)
		}
		for country, countrySlice := range byCountryMap {
			sort.SliceStable(countrySlice, func(i, j int) bool {
				return countrySlice[i].DeliveryTime < countrySlice[j].DeliveryTime
			})
			data[country] = append(data[country], countrySlice[:3])
			data[country] = append(data[country], countrySlice[len(countrySlice)-3:])
		}
		r.cache.Set("email", data, 0)
		res.Email = data
	}
}
