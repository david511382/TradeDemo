package logger

import (
	"fmt"
	errUtil "zerologix-homework/src/pkg/util/error"

	"github.com/rs/zerolog"
)

func Log(name string, msg string, a ...interface{}) {
	if len(a) > 0 {
		msg = fmt.Sprintf(msg, a...)
	}
	LogError(name, errUtil.New(msg, zerolog.InfoLevel))
}

func LogError(name string, err error) {
	if err == nil {
		return
	}

	logger := GetLogger()
	if logger == nil {
		return
	}

	errs := errUtil.Split(err)
	for _, err := range errs {
		logger.Log(name, err)
	}
}
