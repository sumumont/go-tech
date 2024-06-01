package pke

import (
	"go-tech/pkg/apputil"
	"net/http"
	"strings"
)

type PkeConfig struct {
	ModelId int
	Loader  LoadCode
}

type LoadCode func()

var (
	config = PkeConfig{
		ModelId: 1040,
	}
)

func init() {
	LoadConfig(nil)
}

func LoadConfig(pkeConfig *PkeConfig) {
	if pkeConfig != nil {
		config = *pkeConfig
		if config.ModelId == 0 {
			config.ModelId = 1040
		}

	}
	if config.Loader == nil {
		config.Loader = loadCode
	}

	config.Loader()
}

func fill(code int) int {
	return config.ModelId*100000 + code
}

func Short(shortCode int, args ...string) apputil.APIError {
	return apputil.RaiseAPIError(fill(shortCode), args...)
}
func ShortHttp(httpCode int, shortCode int, args ...string) apputil.APIError {
	if httpCode == 0 {
		httpCode = http.StatusOK
	}
	return apputil.RaiseHttpError(httpCode, fill(shortCode), strings.Join(args, " "))
}
func FromRest(httpCode int, restResp *CommResponse) apputil.APIError {
	return &apputil.APIException{
		StatusCode: httpCode,
		Code:       restResp.Code,
		Data:       restResp.Data,
		Msg:        restResp.Msg,
	}
}
