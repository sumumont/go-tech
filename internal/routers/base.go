package routers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go-tech/internal/logging"
	"go-tech/pkg/apputil"
	"go-tech/pkg/pke"
	"math"
	"net/http"
	"time"
)

type RequestProcessor func(c *gin.Context) (interface{}, error)

func wrapper(handler RequestProcessor) func(c *gin.Context) {
	return func(c *gin.Context) {

		data, err := handler(c)
		resp := pke.CommResponse{Data: data}

		statusCode := http.StatusOK
		if err != nil { //process for error
			var h apputil.APIError
			if errors.As(err, &h) {
				statusCode = h.HttpStatus()
				//msg := i18n.TransMsg(h.Error(), i18n.GetLangCtx(c), strconv.Itoa(h.Errno()), h.ErrorData(), i18n.LocalizeAll)
				resp.Code = h.Errno()
				//resp.Msg = msg
			} else { //如果不是规范的错误，则统一返回 "系统内部错误: << .msg >>"
				h = pke.ErrCmSystem
				statusCode = h.HttpStatus()
				msgData := map[string]string{}
				msgData["msg"] = err.Error()
				//msg := i18n.TransMsg(h.Error(), i18n.GetLangCtx(c), strconv.Itoa(h.Errno()), msgData, i18n.LocalizeAll)
				resp.Code = h.Errno()
				//resp.Msg = msg
			}
		}
		c.JSON(statusCode, resp)
	}
}
func loggerHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		start := time.Now()
		c.Next()
		stop := time.Since(start)
		latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000.0))
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		clientUserAgent := c.Request.UserAgent()
		referer := c.Request.Referer()
		dataLength := c.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}
		logging.Info(c).Fields(logrus.Fields{
			"method":     c.Request.Method,
			"path":       path,
			"statusCode": statusCode,
			"cost":       fmt.Sprintf("(%dμs)", latency),
			"clientIP":   clientIP,
			"referer":    referer,
			"dataLength": dataLength,
			"userAgent":  clientUserAgent,
		}).Send()
		if len(c.Errors) > 0 {
			logging.Error(c, errors.New(c.Errors.ByType(gin.ErrorTypePrivate).String())).Send()
		} else {
			if statusCode > 499 {
				logging.Error(c, errors.New("statusCode error")).Msgf("statusCode:%d", statusCode)
			} else if statusCode > 399 {
				logging.Warn(c).Int("statusCode", statusCode).Send()
			} else if statusCode == 200 {
			} else {
				logging.Info(c).Int("statusCode", statusCode).Send()
			}
		}
	}
}
