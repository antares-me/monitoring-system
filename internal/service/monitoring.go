package service

import (
	"context"
	"log"

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
		err        error
	)
	if resultSetT.SMS, err = s.sms.GetResultData(ctx); err != nil {
		log.Println(err)
		status = false
	}
	if resultSetT.MMS, err = s.mms.GetResultData(ctx); err != nil {
		log.Println(err)
		status = false
	}
	if resultSetT.VoiceCall, err = s.voicecall.GetResultData(ctx); err != nil {
		log.Println(err)
		status = false
	}
	if resultSetT.Email, err = s.email.GetResultData(ctx); err != nil {
		log.Println(err)
		status = false
	}
	if resultSetT.Billing, err = s.billing.GetResultData(ctx); err != nil {
		log.Println(err)
		status = false
	}
	if resultSetT.Support, err = s.support.GetResultData(ctx); err != nil {
		log.Println(err)
		status = false
	}
	if resultSetT.Incident, err = s.incident.GetResultData(ctx); err != nil {
		log.Println(err)
		status = false
	}

	return resultSetT, status
}
