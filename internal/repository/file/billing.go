package file

import (
	"context"
	"io/ioutil"
	"math"
	"sync"

	"antares-me/monitoring-system/internal/domain"

	"antares-me/monitoring-system/pkg/cache"
)

type BillingRepo struct {
	file  string
	data  domain.BillingData
	cache *cache.Cache
}

func NewBillingRepo(fp string, c *cache.Cache) *BillingRepo {
	return &BillingRepo{
		file:  fp,
		data:  domain.BillingData{},
		cache: c,
	}
}

func (r *BillingRepo) parseData(ctx context.Context) error {
	var (
		digit uint8
		j     float64
	)
	content, err := ioutil.ReadFile(r.file)
	if err != nil {
		return err
	}
	for i := len(content) - 1; i >= 0; i-- {
		if string(content[i]) == "1" {
			digit += uint8(math.Pow(2, j))
			// Я бы сделал так, но в ТЗ есть возведение в степень.
			// digit |= (1 << j)
		}
		j++
	}
	r.data.CheckoutPage = digit&(1<<5) > 0
	r.data.FraudControl = digit&(1<<4) > 0
	r.data.Recurring = digit&(1<<3) > 0
	r.data.Payout = digit&(1<<2) > 0
	r.data.Purchase = digit&(1<<1) > 0
	r.data.CreateCustomer = digit&(1<<0) > 0
	return nil
}

func (r *BillingRepo) GetResultData(ctx context.Context, wg *sync.WaitGroup, res *domain.ResultSetT, e *[]error) {
	defer wg.Done()
	if val, has := r.cache.Get("billing"); has == true {
		res.Billing = val.(domain.BillingData)
	} else {
		if err := r.parseData(ctx); err != nil {
			*e = append(*e, err)
			return
		}
		r.cache.Set("billing", r.data, 0)
		res.Billing = r.data
	}
}
