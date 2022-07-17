package common

import (
	"zerologix-homework/src/pkg/constant"
	errUtil "zerologix-homework/src/pkg/util/error"
)

func IsContainNotChange(errInfo errUtil.IError) bool {
	if errInfo == nil {
		return false
	}

	for _, err := range errUtil.Split(errInfo) {
		errInfo, ok := err.(errUtil.IError)
		if ok {
			if errUtil.Equal(errInfo, constant.ErrInfoRedisNotChange) {
				return true
			}
			continue
		}

		if err.Error() == constant.ERROR_MSG_REDIS_NOT_CHANGE {
			return true
		}
	}

	return false
}
