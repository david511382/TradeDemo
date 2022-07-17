package timeutil

import (
	"time"
)

type ITime interface {
	Now() time.Time
}
