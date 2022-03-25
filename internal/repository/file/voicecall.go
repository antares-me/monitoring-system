package file

import (
	"context"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"antares-me/monitoring-system/internal/domain"

	"antares-me/monitoring-system/pkg/cache"
	"antares-me/monitoring-system/pkg/country"
	"antares-me/monitoring-system/pkg/provider"
)

type VoiceCallRepo struct {
	file  string
	data  []domain.VoiceCallData
	cache *cache.Cache
}

func NewVoiceCallRepo(fp string, c *cache.Cache) *VoiceCallRepo {
	return &VoiceCallRepo{
		file:  fp,
		data:  []domain.VoiceCallData{},
		cache: c,
	}
}

func (r *VoiceCallRepo) parseData(ctx context.Context) error {
	r.data = []domain.VoiceCallData{}
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
		if delimitersCount != 7 {
			log.Println("VoiceCall: неверное количество полей:", delimitersCount+1)
			continue
		}
		fields := strings.Split(line, ";")
		if validateVoiceCallFields(fields) {
			connectionStability, err := strconv.ParseFloat(fields[4], 32)
			if err != nil {
				continue
			}
			ttfb, err := strconv.Atoi(fields[5])
			if err != nil {
				continue
			}
			voicePurity, err := strconv.Atoi(fields[6])
			if err != nil {
				continue
			}
			medianOfCallsTime, err := strconv.Atoi(fields[7])
			if err != nil {
				continue
			}
			r.data = append(r.data, domain.VoiceCallData{
				Country:             fields[0],
				Bandwidth:           fields[1],
				ResponseTime:        fields[2],
				Provider:            fields[3],
				ConnectionStability: float32(connectionStability),
				TTFB:                ttfb,
				VoicePurity:         voicePurity,
				MedianOfCallsTime:   medianOfCallsTime,
			})
		}
	}
	return nil
}

// validateVoiceCallFields проверяет поля одной VoiceCall записи на соответствие требованиям ТЗ
func validateVoiceCallFields(fields []string) bool {
	if !country.IsRightCode(fields[0]) {
		log.Println("VoiceCall: неверный код страны:", fields[0])
		return false
	}
	if !provider.IsRightCode(fields[3], "VoiceCall") {
		log.Println("VoiceCall: неверный провайдер:", fields[3])
		return false
	}
	return true
}

func (r *VoiceCallRepo) GetResultData(ctx context.Context) ([]domain.VoiceCallData, error) {
	if val, has := r.cache.Get("voicecall"); has == true {
		v := val.([]domain.VoiceCallData)
		return v, nil
	} else {
		if err := r.parseData(ctx); err != nil {
			return []domain.VoiceCallData{}, err
		}
		r.cache.Set("voicecall", r.data, 0)
		return r.data, nil
	}
}
