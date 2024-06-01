package pke

import (
	"go-tech/pkg/apputil"
	"net/http"
)

// CodeCm开头 表示通用的错误码  任何业务系统的意思都一样
const (
	CmSystem         = 0
	CmParam          = 1
	CmIo             = 2
	CmFileOption     = 3
	CmIllegalPath    = 4
	CmJsonmarshal    = 5
	CmJsonUnmarshal  = 6
	CmInvalidSuffix  = 7
	CmInvalidFileExt = 8
	CmUploadFile     = 9
	CmMqOption       = 10
	CmResponseIo     = 11
	CmNotFound       = 12
)
const (
	DbError = 300 + iota
	DbQueryFailed
	DbExecFailed
	DbDuplicate
	DbUpdateUnexpect
	DbWrongType
	DbReadRows
)

const (
	ManifestNotFound    = 105
	ManifestWriteFailed = 106
	ManifestReadFailed  = 107
	ExtNotMd            = 111
	ManifestInvalid     = 112
	EmptyDir            = 208
)
const (
	InvalidModel  = 1001
	ModelNotReady = 1002
)

var (
	ErrCmSystem            apputil.APIError
	ErrCmInvalidParam      apputil.APIError
	ErrCmIo                apputil.APIError
	ErrCmFileOption        apputil.APIError
	ErrCmIllegalPath       apputil.APIError
	ErrCmJsonmarshal       apputil.APIError
	ErrCmJsonUnmarshal     apputil.APIError
	ErrCmInvalidSuffix     apputil.APIError
	ErrCmInvalidFileExt    apputil.APIError
	ErrManifestNotFound    apputil.APIError
	ErrManifestWriteFailed apputil.APIError
	ErrManifestReadFailed  apputil.APIError
	ErrExtNotMd            apputil.APIError
	ErrCmUploadFile        apputil.APIError
	ErrInvalidModel        apputil.APIError
	ErrManifestInvalid     apputil.APIError
	ErrEmptyDir            apputil.APIError
	ErrDbError             apputil.APIError
	ErrDbQueryFailed       apputil.APIError
	ErrDbExecFailed        apputil.APIError
	ErrDbDuplicate         apputil.APIError
	ErrDbUpdateUnexpect    apputil.APIError
	ErrDbWrongType         apputil.APIError
	ErrDbReadRows          apputil.APIError
	ErrCmMqOption          apputil.APIError
	ErrCmResponseIo        apputil.APIError
	ErrModelNotReady       apputil.APIError
	ErrCmNotFound          apputil.APIError
)

func loadCode() {
	ErrCmSystem = Short(CmSystem, "system error")
	ErrCmInvalidParam = Short(CmParam, "invalid params")
	ErrCmIo = Short(CmIo, "io error")
	ErrCmJsonmarshal = Short(CmJsonmarshal, "json marshal failed")
	ErrCmJsonUnmarshal = Short(CmJsonUnmarshal, "json un marshal failed")
	ErrCmInvalidSuffix = Short(CmInvalidSuffix, "file suffix invalid")
	ErrCmFileOption = Short(CmFileOption, "file option failed")
	ErrCmIllegalPath = Short(CmIllegalPath, "illegal path")
	ErrCmInvalidFileExt = Short(CmInvalidFileExt, "invalid image ext")
	ErrCmUploadFile = Short(CmUploadFile, "upload file failed")
	ErrManifestNotFound = Short(ManifestNotFound, "error manifest not found")
	ErrManifestWriteFailed = Short(ManifestWriteFailed, "manifest write failed")
	ErrManifestReadFailed = Short(ManifestReadFailed, "manifest read failed")
	ErrExtNotMd = Short(ExtNotMd, "file ext need .md")
	ErrInvalidModel = Short(InvalidModel, "invalid model")
	ErrManifestInvalid = Short(ManifestInvalid, "manifest invalid")
	ErrEmptyDir = Short(EmptyDir, "error empty dir")
	ErrDbError = Short(DbError, "db error")
	ErrDbQueryFailed = Short(DbQueryFailed, "db query failed")
	ErrDbExecFailed = Short(DbExecFailed, "db sql exec failed")
	ErrDbDuplicate = Short(DbDuplicate, "db data duplication")
	ErrDbUpdateUnexpect = Short(DbUpdateUnexpect, "db update unexpect")
	ErrDbWrongType = Short(DbWrongType, "db wrong type")
	ErrDbReadRows = Short(DbReadRows, "db read rows failed")
	ErrCmMqOption = Short(CmMqOption, "mq option failed")
	ErrCmResponseIo = Short(CmResponseIo, "http response io read error")
	ErrModelNotReady = Short(ModelNotReady, "model not ready")
	ErrCmNotFound = ShortHttp(http.StatusOK, CmNotFound, "resource not found")
}
