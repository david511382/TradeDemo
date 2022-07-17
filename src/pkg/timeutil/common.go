package timeutil

import (
	"zerologix-homework/bootstrap"

	"golang.org/x/sync/singleflight"
)

var (
	t  *Time
	sf singleflight.Group
)

func GetTimeUtil() (ITime, error) {
	if t == nil {
		_, err, _ := sf.Do("key", func() (interface{}, error) {
			cfg, err := bootstrap.Get()
			if err != nil {
				return nil, err
			}
			tu, err := NewTime(cfg)
			if err != nil {
				return nil, err
			}
			t = tu
			return nil, nil
		})
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}
