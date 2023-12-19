package utils

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type EchoUtils struct {
	E echo.Context
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
