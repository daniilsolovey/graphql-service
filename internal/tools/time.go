package tools

import (
	"time"

	"github.com/reconquest/karma-go"
)

var TimeNow = time.Now

func GetCurrentMoscowTime() (time.Time, error) {
	moscowTime, err := GetTimeMoscowLocation()
	if err != nil {
		return time.Time{}, karma.Format(
			err,
			"unable to get moscow location",
		)
	}

	return TimeNow().In(moscowTime), nil

}

func GetTimeMoscowLocation() (*time.Location, error) {
	moscowTime, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to load location: Europe/Moscow",
		)
	}

	return moscowTime, nil
}
