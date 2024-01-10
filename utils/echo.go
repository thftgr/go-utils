package utils

import (
	"github.com/labstack/echo/v4"
	"github.com/thftgr/go-utils/http/header"
	"github.com/thftgr/go-utils/logger"
	"net/http"
)

type EchoUtils struct {
	E echo.Context
}

func GetEchoUtils(c echo.Context) *EchoUtils {
	return &EchoUtils{E: c}
}

func (e *EchoUtils) GetRequestId() (id string) {
	id = e.E.Request().Header.Get(echo.HeaderXRequestID)
	if id == "" {
		id = e.E.Response().Header().Get(echo.HeaderXRequestID)
	}
	return
}

func (e *EchoUtils) GetHeader(key string) (res string) {
	res = e.E.Request().Header.Get(key)
	if res == "" {
		res = e.E.Response().Header().Get(key)
	}
	return
}

func (e *EchoUtils) GetPath() (res string) {
	res = e.E.Request().URL.Path
	if res == "" {
		res = "/"
	}
	return
}

func (e *EchoUtils) GetStatus() (statusCode int, status string) {
	statusCode = e.E.Response().Status
	status = http.StatusText(statusCode)
	return
}

func (e *EchoUtils) GetLogger() logger.GroupLogger {
	switch v := e.E.Get("logger").(type) {
	case logger.GroupLogger:
		return v
	default:
		return nil
	}
}

func (e *EchoUtils) SetLogger(l logger.GroupLogger) {
	e.E.Set("logger", l)
}

func (e *EchoUtils) Accept() (accept []header.Accept) {
	return header.ParseAccept(e.GetHeader(echo.HeaderAccept))
}

func (e *EchoUtils) ContentType() (contentType *header.ContentType) {
	return header.ParseContentType(e.GetHeader(echo.HeaderContentType))
}
