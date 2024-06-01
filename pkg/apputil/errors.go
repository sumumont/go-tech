package apputil

import (
	"net/http"
	"strings"
)

var APP_UTIL_MODULE_ID = 0

const MAX_ERROR_CODE = 100000

const (
	ERROR_SUCCESS             = 0
	ERROR_NOT_FOUND           = 1
	ERROR_PARAM_ERROR         = 2
	ERROR_ALREADY_EXISTS      = 3
	ERROR_NO_API              = 4
	ERROR_NOT_IMPLEMENT       = 5
	ERROR_UNKNOWN_ERROR       = 6
	ERROR_FILE_NOT_FOUND      = 7
	ERROR_PATH_NOT_FOUND      = 8
	ERROR_OS_IO_ERROR         = 9
	ERROR_OS_READ_DIR         = 10
	ERROR_OS_REMOVE_FILE      = 11
	ERROR_OS_OPEN_FILE        = 12
	ERROR_FILE_TOO_LARGE      = 13
	ERROR_FILE_TYPE_ERROR     = 14
	ERROR_FILE_ALREADY_EXISTS = 15
	ERROR_FILE_PERMISSION     = 16
	ERROR_OS_RENAME_FILE      = 17
	ERROR_OS_CREATE_DIR       = 18
	ERROR_OS_STAT_FILE        = 19

	ERROR_SERVER_BUSY  = 22
	ERROR_WOULD_BLOCK  = 23
	ERROR_NO_AUTH      = 24
	ERROR_STILL_ACTIVE = 25
	ERROR_LOGIC_ERROR  = 26
	//ERROR_SINGLETON_RUN_EXISTS = 27
	//ERROR_RUN_CANNOT_RESTART   = 28
	ERROR_REMOTE_NETWORK_ERROR = 29
	ERROR_REMOTE_REST_ERROR    = 30
	ERROR_REMOTE_GRPC_ERROR    = 31
	ERROR_CANNOT_COMMIT        = 32
	ERROR_EXCEED_LIMIT         = 33
	//ERROR_INVALID_ENDPOINT_STATUS = 34
	ERROR_INVALID_FORMAT = 35
	ERROR_MQ_SEND_ERROR  = 36

	ERROR_PROXY_IO_ERROR         = 37
	ERROR_PROXY_INVALID_RESPONSE = 38

	ERROR_UNKNOWN_EXTEND_ANY_PROTO = 39
	ERROR_USER_CANCEL              = 40
	ERROR_JSON_UNMARSHAL           = 41
	ERROR_JSON_MARSHAL             = 42
	ERROR_USER_PAUSED              = 43
	ERROR_INVALID_STATUS           = 44
	ERROR_USER_DISCARD             = 45

	ERROR_PERMISSION_ERROR = 60
)
const (
	ERROR_DB_BEGIN            = 300
	ERROR_DB_QUERY_FAILED     = 301
	ERROR_DB_EXEC_FAILED      = 302
	ERROR_DB_DUPLICATE        = 303
	ERROR_DB_UPDATE_UNEXPECT  = 304
	ERROR_DB_WRONG_TYPE       = 305
	ERROR_DB_READ_ROWS        = 306
	ERROR_DB_VERSION_MISMATCH = 307
	ERROR_DB_WRONG_SCHEMA     = 308
	ERROR_DB_INIT_FAILED      = 309
	ERROR_DB_FLUSH_FAILED     = 310
)

type APIError interface {
	Error() string
	Errno() int
	HttpStatus() int
	ErrorData() interface{}

	WithData(data interface{}) APIError
	WithMsg(msg string) APIError
	WithHttpStatus(httpCode int) APIError
	//@add:  [nil safe] check no error
	NoError() bool
}

// @add: support multiple errors aggregate
type APIErrors interface {
	APIError
	AddError(err APIError)
	//@mark: [nil safe] get all errors slice
	GetErrors() []APIError
	//@mark: [nil safe] if zero errors will return nil
	ToError() APIErrors
}

type APIException struct {
	StatusCode int    `json:"-"` // http status code ,should be rarely used
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
	//@add: support for i18n template error message format
	Data interface{} `json:"data,omitempty"`
}

func (e *APIException) Error() string {
	return e.Msg
}
func (e *APIException) Errno() int {
	return e.Code
}
func (e *APIException) ErrorData() interface{} {
	return e.Data
}
func (e *APIException) HttpStatus() int {
	return e.StatusCode
}

func (e *APIException) WithData(data interface{}) APIError {
	e.Data = data
	return e
}
func (e *APIException) WithMsg(msg string) APIError {
	e.Msg = msg
	return e
}
func (e *APIException) WithHttpStatus(httpCode int) APIError {
	e.StatusCode = httpCode
	return e
}
func (e *APIException) NoError() bool {
	return e == nil || e.Code == 0
}

func NewAPIError(code int) APIError {
	return &APIException{
		Code:       code,
		StatusCode: http.StatusBadRequest,
	}
}
func NewAPIErrorWith(code int, msg string, data interface{}) APIError {
	return &APIException{
		Code:       code,
		Msg:        msg,
		Data:       data,
		StatusCode: http.StatusBadRequest,
	}
}

// @add: support multiple errors management
type multiAPIExceptions struct {
	APIException
}

func (e *multiAPIExceptions) AddError(err APIError) {
	if err == nil {
		return
	}
	errors, _ := e.Data.([]APIError)
	if multi_err, _ := err.(APIErrors); multi_err != nil {
		e.Data = append(errors, multi_err.GetErrors()...)
	} else if err != nil {
		e.Data = append(errors, err)
	}
}

// @mark: safe transform to APIError normal interface
func (e *multiAPIExceptions) ToError() APIErrors {
	if e.NoError() {
		return nil
	}
	return e
}
func (e *multiAPIExceptions) GetErrors() []APIError {
	if e == nil {
		return nil
	}
	errors, _ := e.Data.([]APIError)
	return errors
}
func (e *multiAPIExceptions) WithData(data interface{}) APIError {
	panic("APIErrors cannot set user data !!!")
}
func (e *multiAPIExceptions) NoError() bool {
	return e == nil || len(e.GetErrors()) == 0
}

func NewAPIErrors(code int, msg ...string) APIErrors {
	return &multiAPIExceptions{
		APIException: APIException{
			StatusCode: http.StatusBadRequest,
			Code:       code,
			Msg:        strings.Join(msg, " "),
		},
	}
}

func NewAPIException(statusCode, code int, msg string) *APIException {
	return &APIException{
		StatusCode: statusCode,
		Code:       code,
		Msg:        msg,
	}
}

func UnAuthorizedError(msg string) *APIException {
	return NewAPIException(http.StatusUnauthorized, ERROR_NO_AUTH, msg)
}

func NotFoundError() *APIException {
	return NewAPIException(http.StatusNotFound, ERROR_NOT_FOUND, http.StatusText(http.StatusNotFound))
}

func UnknownError(msg string) *APIException {
	return NewAPIException(http.StatusForbidden, ERROR_UNKNOWN_ERROR, msg)
}

func PermissionError(msg string) *APIException {
	return NewAPIException(http.StatusForbidden, ERROR_PERMISSION_ERROR, msg)
}

func ParameterError(msg string) *APIException {
	return NewAPIException(http.StatusBadRequest, ERROR_PARAM_ERROR, msg)
}

func NotImplementError(msg string) *APIException {
	return NewAPIException(http.StatusBadRequest, ERROR_NOT_IMPLEMENT, msg)
}

func ReqWouldBlockError(msg string) *APIException {
	return NewAPIException(http.StatusOK, ERROR_WOULD_BLOCK, msg)
}

func ReqUserCancel(msg string) *APIException {
	return NewAPIException(http.StatusInternalServerError, ERROR_USER_CANCEL, msg)
}

func LogicError(msg string) *APIException {
	return NewAPIException(http.StatusInternalServerError, ERROR_LOGIC_ERROR, msg)
}

// server error
func RaiseServerError(code int, args ...string) APIError {
	return NewAPIException(http.StatusInternalServerError, code, strings.Join(args, " "))
}

// client error
func RaiseAPIError(code int, args ...string) APIError {
	return NewAPIException(http.StatusBadRequest, code, strings.Join(args, " "))
}

func RaiseHttpError(statusCode int, code int, status string) APIError {
	return NewAPIException(statusCode, code, status)
}

func CheckWithApiError(err error, code int) APIError {
	if err != nil {
		if err1, ok := err.(APIError); ok {
			return err1
		}
		return RaiseAPIError(code, err.Error())
	} else if code != 0 {
		return RaiseAPIError(code)
	} else {
		return nil
	}
}

// probe some usually occur DB error , compatible with old xxxyyyyy code formats !
func IsDBErrorNotFound(err APIError) bool {
	return err.Errno()%MAX_ERROR_CODE == ERROR_NOT_FOUND
}
func IsDBErrorDuplicate(err APIError) bool {
	return err != nil && err.Errno()%MAX_ERROR_CODE == ERROR_DB_DUPLICATE
}
func IsDBErrorUnexpectUpdate(err APIError) bool {
	return err.Errno()%MAX_ERROR_CODE == ERROR_DB_UPDATE_UNEXPECT
}

// @add: probe whether api call result is completed
func IsErrorWouldBlock(err APIError) bool {
	return err != nil && err.Errno() == ERROR_WOULD_BLOCK
}

func SetModuleInfo(moduleId int) {
	if APP_UTIL_MODULE_ID != 0 && APP_UTIL_MODULE_ID != moduleId {
		panic("APP_UTIL_MODULE_ID has been set another value already !!!")
	}
	APP_UTIL_MODULE_ID = moduleId
}
