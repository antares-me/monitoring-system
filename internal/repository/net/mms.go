package net

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"

	"antares-me/monitoring-system/internal/domain"

	"antares-me/monitoring-system/pkg/cache"
	"antares-me/monitoring-system/pkg/country"
	"antares-me/monitoring-system/pkg/provider"
)

type MmsRepo struct {
	url   string
	data  []domain.MMSData
	cache *cache.Cache
}

func NewMmsRepo(u string, c *cache.Cache) *MmsRepo {
	return &MmsRepo{
		url:   u,
		data:  []domain.MMSData{},
		cache: c,
	}
}

func (r *MmsRepo) parseData(ctx context.Context) error {
	r.data = []domain.MMSData{}
	var d []domain.MMSData
	resp, err := http.Get(r.url)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("MMS: неверный код ответа сервера: %d", resp.StatusCode)
	}
	content, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, &d)
	if err != nil {
		return err
	}
	for _, dStruct := range d {
		if validateMmsFields(dStruct) {
			r.data = append(r.data, dStruct)
		}
	}
	return nil
}

// validateMmsFields проверяет поля одной ММС записи на соответствие требованиям ТЗ
func validateMmsFields(fields domain.MMSData) bool {
	if !country.IsRightCode(fields.Country) {
		log.Println("MMS: неверный код страны:", fields.Country)
		return false
	}
	if !provider.IsRightCode(fields.Provider, "Mms") {
		log.Println("MMS: неверный провайдер:", fields.Provider)
		return false
	}
	return true
}

func (r *MmsRepo) replaceCountryCodes() {
	for i, mms := range r.data {
		r.data[i].Country = country.ReplaceCountruCode(mms.Country)
	}
}

// SortByProvider сортирует срез данных MMS по провайдеру
func (r *MmsRepo) sortByProvider() []domain.MMSData {
	data := make([]domain.MMSData, len(r.data))
	copy(data, r.data)
	sort.SliceStable(data, func(i, j int) bool {
		return data[i].Provider < data[j].Provider
	})
	return data
}

// SortByProvider сортирует срез данных MMS по стране
func (r *MmsRepo) sortByCountry() []domain.MMSData {
	data := make([]domain.MMSData, len(r.data))
	copy(data, r.data)
	sort.SliceStable(data, func(i, j int) bool {
		return data[i].Country < data[j].Country
	})
	return data
}

func (r *MmsRepo) GetResultData(ctx context.Context) ([][]domain.MMSData, error) {
	if val, has := r.cache.Get("mms"); has == true {
		v := val.([][]domain.MMSData)
		return v, nil
	} else {
		data := [][]domain.MMSData{}
		if err := r.parseData(ctx); err != nil {
			return data, err
		}
		r.replaceCountryCodes()
		data = append(data, r.sortByProvider())
		data = append(data, r.sortByCountry())
		r.cache.Set("mms", data, 0)
		return data, nil
	}
}
