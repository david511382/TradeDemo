package errorcode

import (
	errUtil "zerologix-homework/src/pkg/util/error"
)

type IErrorcode interface {
	Log() (code int, err error)
	error
}

type Errorcode struct {
	error
	logErrInfo error
	code       int
}

func newErrorcode(msg ErrorMsg, logErrInfo error, code int) IErrorcode {
	err := errUtil.NewRaw(string(msg))
	return Errorcode{
		error:      err,
		logErrInfo: logErrInfo,
		code:       code,
	}
}

func (ec Errorcode) Log() (code int, err error) {
	err = ec.logErrInfo
	code = ec.code
	return
}

func IsContain(err error, errMsg ErrorMsg) bool {
	errCode := GetErrorcode(err)
	if errCode == nil {
		return false
	}
	return errCode.Error() == errMsg.Error()
}

func GetErrorcode(err error) IErrorcode {
	if err == nil {
		return nil
	}
	for _, err := range errUtil.Split(err) {
		errCode, ok := err.(IErrorcode)
		if ok {
			return errCode
		}
	}
	return nil
}
