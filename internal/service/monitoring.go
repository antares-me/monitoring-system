package service

import (
	"context"
	"log"
	"sync"

	"antares-me/monitoring-system/internal/domain"
	"antares-me/monitoring-system/internal/repository"
)

type MonitoringService struct {
	sms       repository.Sms
	mms       repository.Mms
	voicecall repository.VoiceCall
	email     repository.Email
	billing   repository.Billing
	incident  repository.Incident
	support   repository.Support
}

func NewMonitoringService(r *repository.Repositories) *MonitoringService {
	return &MonitoringService{
		sms:       r.Sms,
		mms:       r.Mms,
		voicecall: r.VoiceCall,
		email:     r.Email,
		billing:   r.Billing,
		incident:  r.Incident,
		support:   r.Support,
	}
}

func (s *MonitoringService) GetStatus(ctx context.Context) domain.ResultT {
	data := domain.ResultT{}
	resultSet, status := s.getResultData(ctx)
	data.Status = status
	if !status {
		data.Error = "Error on collect data"
	} else {
		data.Data = resultSet
	}
	return data
}

func (s *MonitoringService) getResultData(ctx context.Context) (domain.ResultSetT, bool) {
	var (
		status     bool = true
		resultSetT domain.ResultSetT
		err        []error
	)

	wg := &sync.WaitGroup{}
	wg.Add(7)
	go s.sms.GetResultData(ctx, wg, &resultSetT, &err)
	go s.mms.GetResultData(ctx, wg, &resultSetT, &err)
	go s.voicecall.GetResultData(ctx, wg, &resultSetT, &err)
	go s.email.GetResultData(ctx, wg, &resultSetT, &err)
	go s.billing.GetResultData(ctx, wg, &resultSetT, &err)
	go s.support.GetResultData(ctx, wg, &resultSetT, &err)
	go s.incident.GetResultData(ctx, wg, &resultSetT, &err)
	wg.Wait()
	if len(err) > 0 {
		for _, e := range err {
			log.Printf("Error: %s\n", e)
		}
		status = false
	}
	return resultSetT, status
}
