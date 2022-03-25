package net

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"

	"antares-me/monitoring-system/internal/domain"

	"antares-me/monitoring-system/pkg/cache"
)

type IncidentRepo struct {
	url   string
	data  []domain.IncidentData
	cache *cache.Cache
}

func NewIncidentRepo(u string, c *cache.Cache) *IncidentRepo {
	return &IncidentRepo{
		url:   u,
		data:  []domain.IncidentData{},
		cache: c,
	}
}

func (r *IncidentRepo) parseData(ctx context.Context) error {
	resp, err := http.Get(r.url)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Incident: неверный код ответа сервера: %d", resp.StatusCode)
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if err := json.Unmarshal(content, &r.data); err != nil {
		return err
	}
	return nil
}

func (r *IncidentRepo) GetResultData(ctx context.Context) ([]domain.IncidentData, error) {
	if val, has := r.cache.Get("incident"); has == true {
		v := val.([]domain.IncidentData)
		return v, nil
	} else {
		if err := r.parseData(ctx); err != nil {
			return r.data, err
		}
		sort.SliceStable(r.data, func(i, j int) bool {
			return r.data[i].Status < r.data[j].Status
		})
		r.cache.Set("incident", r.data, 0)
		return r.data, nil
	}
}
