package service

import (
	model "commentservice/internal/model/domain"
	"time"
)

type TimeService struct {
}

func (t *TimeService) GetTimeNowUtc() time.Time {
	timeNowString := time.Now().UTC().Format(model.TimeLayout)
	timeNow, _ := time.Parse(model.TimeLayout, timeNowString)

	return timeNow
}
