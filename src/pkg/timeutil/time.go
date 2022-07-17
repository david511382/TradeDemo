package timeutil

import (
	"time"
	"zerologix-homework/bootstrap"
)

type Time struct {
	location *time.Location
}

func NewTime(cfg *bootstrap.Config) (*Time, error) {
	var loc *time.Location
	if cfg.Var.TimeZone == "" {
		loc = time.Local
	} else {
		l, err := time.LoadLocation(cfg.Var.TimeZone)
		if err != nil {
			return nil, err
		}
		loc = l
	}

	return &Time{
		location: loc,
	}, nil
}

func (t Time) Now() time.Time {
	return time.Now().In(t.location)
}
