package middlewares

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/thftgr/go-utils/http/header"
	"net/http"
	"testing"
	"time"
)

const _TEST_PORT = "8080"
const _TEST_URL = "http://127.0.0.1:" + _TEST_PORT

func TestCheckAllowedAccept(t *testing.T) {
	e := echo.New()
	e.HideBanner = true
	e.Any("/application/json",
		func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]interface{}{"message": "ok"})
		},
		CheckAllowedAccept("application/json"),
	)
	e.Any("/application/all",
		func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]interface{}{"message": "ok"})
		},
		CheckAllowedAccept("application/json", "application/xml"),
	)
	shudown := func() {
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
		defer cancel()
		fmt.Println(e.Shutdown(ctx))
		fmt.Println("----------------------------shutdown----------------------------")
	}
	go func() {
		if err := e.Start(":" + _TEST_PORT); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatalf("shutting down the server. err: %+v", err)
		}
	}()
	time.Sleep(time.Second)
	defer shudown()
	//=================================================================================================================

	tests := []struct {
		name       string
		url        string
		accept     header.MimeType
		statusCode int // statusCode
	}{
		{name: "", url: "/application/json", accept: "", statusCode: http.StatusOK},
		{name: "", url: "/application/json", accept: "*/*", statusCode: http.StatusOK},

		{name: "", url: "/application/json", accept: "application/json", statusCode: http.StatusOK},
		{name: "", url: "/application/json", accept: "application/xml", statusCode: http.StatusNotAcceptable},

		{name: "", url: "/application/json", accept: "application/xml,*/*", statusCode: http.StatusOK},
		{name: "", url: "/application/json", accept: "application/xml,application/json", statusCode: http.StatusOK},

		{name: "", url: "/application/all", accept: "application/xml", statusCode: http.StatusOK},
		{name: "", url: "/application/all", accept: "application/json", statusCode: http.StatusOK},
		{name: "", url: "/application/all", accept: "text/html", statusCode: http.StatusNotAcceptable},

		{name: "", url: "/application/all", accept: "*/json", statusCode: http.StatusNotAcceptable},
		{name: "", url: "/application/all", accept: "/json", statusCode: http.StatusNotAcceptable},
		{name: "", url: "/application/all", accept: "application/*", statusCode: http.StatusOK},
		{name: "", url: "/application/all", accept: "application/", statusCode: http.StatusNotAcceptable},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			request, err := http.NewRequest("GET", _TEST_URL+tt.url, nil)
			if err != nil {
				t.Error(err)
				return
			}
			request.Header.Set("accept", tt.accept.String())
			client := http.Client{Timeout: time.Second}
			res, err := client.Do(request)
			if err != nil {
				t.Error(err)
				return
			}
			defer res.Body.Close()
			if res.StatusCode != tt.statusCode {
				t.Errorf("CheckAllowedAccept() = %v, statusCode %v", res.StatusCode, tt.statusCode)
			}
		})
	}
}

func TestCheckAllowedContentType(t *testing.T) {
	e := echo.New()
	e.HideBanner = true
	e.Any("/application/json",
		func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]interface{}{"message": "ok"})
		},
		CheckAllowedContentType("application/json"),
	)
	e.Any("/application/all",
		func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]interface{}{"message": "ok"})
		},
		CheckAllowedContentType("application/json", "application/xml"),
	)
	shudown := func() {
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
		defer cancel()
		fmt.Println(e.Shutdown(ctx))
		fmt.Println("----------------------------shutdown----------------------------")
	}
	go func() {
		if err := e.Start(":" + _TEST_PORT); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatalf("shutting down the server. err: %+v", err)
		}
	}()
	time.Sleep(time.Second)
	defer shudown()
	//=================================================================================================================

	tests := []struct {
		name        string
		url         string
		contentType header.MimeType
		statusCode  int // statusCode
	}{
		{name: "", url: "/application/json", contentType: "*/*", statusCode: http.StatusUnsupportedMediaType},
		{name: "", url: "/application/json", contentType: "application/json", statusCode: http.StatusOK},
		{name: "", url: "/application/json", contentType: "application/xml", statusCode: http.StatusUnsupportedMediaType},

		{name: "", url: "/application/all", contentType: "application/xml", statusCode: http.StatusOK},
		{name: "", url: "/application/all", contentType: "application/json", statusCode: http.StatusOK},
		{name: "", url: "/application/all", contentType: "text/html", statusCode: http.StatusUnsupportedMediaType},

		{name: "", url: "/application/all", contentType: "*/json", statusCode: http.StatusUnsupportedMediaType},
		{name: "", url: "/application/all", contentType: "/json", statusCode: http.StatusUnsupportedMediaType},
		{name: "", url: "/application/all", contentType: "application/*", statusCode: http.StatusUnsupportedMediaType},
		{name: "", url: "/application/all", contentType: "application/", statusCode: http.StatusUnsupportedMediaType},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			request, err := http.NewRequest("GET", _TEST_URL+tt.url, nil)
			if err != nil {
				t.Error(err)
				return
			}
			request.Header.Set("Content-Type", tt.contentType.String())
			client := http.Client{Timeout: time.Second}
			res, err := client.Do(request)
			if err != nil {
				t.Error(err)
				return
			}
			defer res.Body.Close()
			if res.StatusCode != tt.statusCode {
				t.Errorf("CheckAllowedAccept() = %v, statusCode %v", res.StatusCode, tt.statusCode)
			}
		})
	}
}
