package common

import (
	"encoding/json"
	"zerologix-homework/src/pkg/errorcode"
	"zerologix-homework/src/server/domain"
	"zerologix-homework/src/server/resp"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data any) {
	response(
		c,
		errorcode.ERROR_MSG_SUCCESS.New(nil),
		&resp.Base{
			Data: data,
		},
	)
}

func FailRequest(c *gin.Context, err error) {
	FailErrorcode(c, errorcode.ERROR_MSG_REQUEST.New(err))
}

func FailAuth(c *gin.Context, err error) {
	FailErrorcode(c, errorcode.ERROR_MSG_AUTH.New(err))
}

func FailForbidden(c *gin.Context, err error) {
	FailErrorcode(c, errorcode.ERROR_MSG_FORBIDDEN.New(err))
}

// err 不得為 nil
func Fail(c *gin.Context, err error) {
	errCode := errorcode.GetErrorcode(err)
	if errCode != nil {
		FailErrorcode(c, errCode)
	} else {
		Fail(c, errorcode.ERROR_MSG_ERROR.New(err))
	}
}

// errCode 不得為 nil
func FailErrorcode(c *gin.Context, errCode errorcode.IErrorcode) {
	response(c, errCode, nil)
}

// errCode 不得為 nil
func response(c *gin.Context, errCode errorcode.IErrorcode, result *resp.Base) {
	if result == nil {
		result = &resp.Base{}
	}
	result.Message = errCode.Error()
	if bs, err := json.Marshal(result); err == nil {
		c.Set(domain.KEY_RESPONSE_CONTEXT, string(bs))
	}

	code, err := errCode.Log()
	if err != nil {
		c.Set(domain.KEY_RESPONSE_ERROR, err)
	}

	c.AbortWithStatusJSON(code, result)
}
