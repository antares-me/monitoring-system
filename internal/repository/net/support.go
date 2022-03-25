package net

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"

	"antares-me/monitoring-system/internal/domain"

	"antares-me/monitoring-system/pkg/cache"
)

type SupportRepo struct {
	url   string
	data  []domain.SupportData
	cache *cache.Cache
}

func NewSupportRepo(u string, c *cache.Cache) *SupportRepo {
	return &SupportRepo{
		url:   u,
		data:  []domain.SupportData{},
		cache: c,
	}
}

func (r *SupportRepo) parseData(ctx context.Context) error {
	resp, err := http.Get(r.url)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Support: неверный код ответа сервера: %d", resp.StatusCode)
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

func (r *SupportRepo) GetResultData(ctx context.Context) ([]int, error) {
	if val, has := r.cache.Get("support"); has == true {
		v := val.([]int)
		return v, nil
	} else {
		timePerTicket := math.Round(60 / 18)
		ticketSumm := 0
		data := make([]int, 2)
		if err := r.parseData(ctx); err != nil {
			return data, err
		}
		for _, v := range r.data {
			ticketSumm += v.ActiveTickets
		}
		data[1] = ticketSumm * int(timePerTicket)
		switch {
		case data[1] < 9:
			data[0] = 1
		case data[1] > 16:
			data[0] = 3
		default:
			data[0] = 2
		}
		r.cache.Set("support", data, 0)
		return data, nil
	}
}
