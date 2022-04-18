package repository

import (
	"context"
	"sync"

	"antares-me/monitoring-system/internal/config"
	"antares-me/monitoring-system/internal/domain"
	"antares-me/monitoring-system/internal/repository/file"
	"antares-me/monitoring-system/internal/repository/net"

	"antares-me/monitoring-system/pkg/cache"
)

type Sms interface {
	GetResultData(ctx context.Context, wg *sync.WaitGroup, res *domain.ResultSetT, err *[]error)
}

type Mms interface {
	GetResultData(ctx context.Context, wg *sync.WaitGroup, res *domain.ResultSetT, err *[]error)
}

type VoiceCall interface {
	GetResultData(ctx context.Context, wg *sync.WaitGroup, res *domain.ResultSetT, err *[]error)
}

type Email interface {
	GetResultData(ctx context.Context, wg *sync.WaitGroup, res *domain.ResultSetT, err *[]error)
}

type Incident interface {
	GetResultData(ctx context.Context, wg *sync.WaitGroup, res *domain.ResultSetT, err *[]error)
}

type Billing interface {
	GetResultData(ctx context.Context, wg *sync.WaitGroup, res *domain.ResultSetT, err *[]error)
}

type Support interface {
	GetResultData(ctx context.Context, wg *sync.WaitGroup, res *domain.ResultSetT, err *[]error)
}

type Repositories struct {
	Sms       Sms
	Mms       Mms
	VoiceCall VoiceCall
	Email     Email
	Incident  Incident
	Billing   Billing
	Support   Support
}

func NewRepositories(cfg *config.Config, c *cache.Cache) *Repositories {
	return &Repositories{
		Sms:       file.NewSmsRepo(cfg.DataFilePath.Sms, c),
		Mms:       net.NewMmsRepo(cfg.DataUrl.Mms, c),
		VoiceCall: file.NewVoiceCallRepo(cfg.DataFilePath.VoiceCall, c),
		Email:     file.NewEmailRepo(cfg.DataFilePath.Email, c),
		Incident:  net.NewIncidentRepo(cfg.DataUrl.Incident, c),
		Billing:   file.NewBillingRepo(cfg.DataFilePath.Billing, c),
		Support:   net.NewSupportRepo(cfg.DataUrl.Support, c),
	}
}
