package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"
	"zerologix-homework/src/logger"
	errUtil "zerologix-homework/src/pkg/util/error"
	"zerologix-homework/src/server/common"
	"zerologix-homework/src/server/domain"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func Logger() gin.HandlerFunc {
	notlogged := []string{
		"/",
	}

	var skip map[string]struct{}
	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		resultErr := errUtil.New(
			"API LOG",
			zerolog.InfoLevel,
		)
		if bs, err := ioutil.ReadAll(c.Request.Body); err == nil {
			resultErr.Attr("Body", string(bs))
			c.Request.Body = io.NopCloser(bytes.NewReader(bs))
		}

		// Process request
		c.Next()

		// Stop timer
		nowTime := time.Now()

		if _, ok := skip[path]; ok {
			return
		}

		// Log only when path is not being skipped

		resultErr.Attr("ClientIP", c.ClientIP())
		resultErr.Attr("Method", c.Request.Method)
		if raw != "" {
			path = path + "?" + raw
		}
		resultErr.Attr("Path", path)
		resultErr.Attr("Proto", c.Request.Proto)
		status := c.Writer.Status()
		resultErr.Attr("Status", status)
		if status == http.StatusInternalServerError {
			resultErr.SetLevel(zerolog.WarnLevel)
		}
		resultErr.Attr("Duration", nowTime.Sub(start).String())
		resultErr.Attr("UserAgent", c.Request.UserAgent())
		if responseValue, isExist := c.Get(domain.KEY_RESPONSE_CONTEXT); isExist {
			resultErr.Attr("Response", responseValue)
		}
		if claims := common.GetClaims(c); claims != nil {
			bs, err := json.Marshal(claims)
			if err == nil {
				resultErr.Attr("Claims", string(bs))
			}
		}
		if value, isExist := c.Get(domain.KEY_RESPONSE_ERROR); isExist {
			if err, ok := value.(error); ok {
				resultErr.Attr("Error", err.Error())
			}
		}

		logger.LogError(logger.NAME_API, resultErr)
	}
}
