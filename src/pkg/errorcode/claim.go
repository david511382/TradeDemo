package errorcode

import (
	"net/http"
)

type ErrorMsg string

const (
	ERROR_MSG_SUCCESS          ErrorMsg = "完成"
	ERROR_MSG_ERROR            ErrorMsg = "發生錯誤"
	ERROR_MSG_TRANSACTION_FAIL ErrorMsg = "交易失敗"

	ERROR_MSG_REQUEST ErrorMsg = "參數錯誤"

	ERROR_MSG_AUTH ErrorMsg = "未登入"

	ERROR_MSG_FORBIDDEN ErrorMsg = "沒權限"
)

func (em ErrorMsg) New(logErrInfos ...error) IErrorcode {
	var err error
	if len(logErrInfos) > 0 {
		err = logErrInfos[0]
	}

	code := http.StatusOK
	switch em {
	case ERROR_MSG_REQUEST:
		code = http.StatusBadRequest
	case ERROR_MSG_AUTH:
		code = http.StatusUnauthorized
	case ERROR_MSG_FORBIDDEN:
		code = http.StatusForbidden
	}

	return newErrorcode(em, err, code)
}

func (em ErrorMsg) Error() string {
	return string(em)
}
