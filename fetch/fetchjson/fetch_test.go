package fetchjson

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"testing"
	"time"
)

const _TEST_PORT = "8080"
const _TEST_URL = "http://127.0.0.1:" + _TEST_PORT

func runEcho() (shudown func()) {
	e := echo.New()
	e.HideBanner = true
	e.Any("/get/200", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{"message": "ok"})
	})
	e.Any("/get/307", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/get/200")
	})
	e.Any("/get/400", func(c echo.Context) error {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "failed"})
	})
	e.Any("/get/404", func(c echo.Context) error {
		return c.JSON(http.StatusNotFound, map[string]interface{}{"message": "failed"})
	})
	e.Any("/get/500", func(c echo.Context) error {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "failed"})
	})
	shudown = func() {
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
	return
}
func TestJsonFetch_GET(t *testing.T) {
	shutdown := runEcho()
	defer shutdown()

	time.Sleep(time.Second)
	fetch := NewJsonFetch(nil, time.Second*5)

	tests := []struct {
		name             string
		fetch            *JsonFetch
		url              string
		wantResponseBody []byte
		wantErr          bool
	}{
		{"TestJsonFetch_GET_001", fetch, _TEST_URL + "/get/200", []byte(`{"message":"ok"}`), false},
		{"TestJsonFetch_GET_002", fetch, _TEST_URL + "/get/307", []byte(`{"message":"ok"}`), false},
		{"TestJsonFetch_GET_003", fetch, _TEST_URL + "/get/400", []byte(`{"message":"failed"}`), true},
		{"TestJsonFetch_GET_004", fetch, _TEST_URL + "/get/404", []byte(`{"message":"failed"}`), true},
		{"TestJsonFetch_GET_005", fetch, _TEST_URL + "/get/500", []byte(`{"message":"failed"}`), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.fetch
			gotResponseBody, err := r.Get(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("GET() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(gotResponseBody) == string(tt.wantResponseBody) {
				t.Errorf("GET() gotResponseBody = %v, want %v", string(gotResponseBody), string(tt.wantResponseBody))
			} else {
				t.Logf("GET() gotResponseBody = %v, want %v", string(gotResponseBody), string(tt.wantResponseBody))
			}
		})
	}
}
func TestJsonFetch_HEAD(t *testing.T) {
	shutdown := runEcho()
	defer shutdown()

	time.Sleep(time.Second)
	fetch := NewJsonFetch(nil, time.Second*5)

	tests := []struct {
		name             string
		fetch            *JsonFetch
		url              string
		wantResponseBody []byte
		wantErr          bool
	}{
		{"TestJsonFetch_HEAD_001", fetch, _TEST_URL + "/get/200", []byte(`{"message":"ok"}`), false},
		{"TestJsonFetch_HEAD_002", fetch, _TEST_URL + "/get/307", []byte(`{"message":"ok"}`), false},
		{"TestJsonFetch_HEAD_003", fetch, _TEST_URL + "/get/400", []byte(`{"message":"failed"}`), true},
		{"TestJsonFetch_HEAD_004", fetch, _TEST_URL + "/get/404", []byte(`{"message":"failed"}`), true},
		{"TestJsonFetch_HEAD_005", fetch, _TEST_URL + "/get/500", []byte(`{"message":"failed"}`), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.fetch
			gotResponseBody, err := r.Head(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("HEAD() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(gotResponseBody) == string(tt.wantResponseBody) {
				t.Errorf("HEAD() gotResponseBody = %v, want %v", string(gotResponseBody), string(tt.wantResponseBody))
			} else {
				t.Logf("HEAD() gotResponseBody = %v, want %v", string(gotResponseBody), string(tt.wantResponseBody))
			}
		})
	}
}
func TestJsonFetch_POST(t *testing.T) {
	shutdown := runEcho()
	defer shutdown()

	time.Sleep(time.Second)
	fetch := NewJsonFetch(nil, time.Second*5)

	tests := []struct {
		name             string
		fetch            *JsonFetch
		url              string
		wantResponseBody []byte
		wantErr          bool
	}{
		{"TestJsonFetch_POST_001", fetch, _TEST_URL + "/get/200", []byte(`{"message":"ok"}`), false},
		{"TestJsonFetch_POST_002", fetch, _TEST_URL + "/get/307", []byte(`{"message":"ok"}`), false},
		{"TestJsonFetch_POST_003", fetch, _TEST_URL + "/get/400", []byte(`{"message":"failed"}`), true},
		{"TestJsonFetch_POST_004", fetch, _TEST_URL + "/get/404", []byte(`{"message":"failed"}`), true},
		{"TestJsonFetch_POST_005", fetch, _TEST_URL + "/get/500", []byte(`{"message":"failed"}`), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.fetch
			gotResponseBody, err := r.Post(tt.url, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("POST() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(gotResponseBody) == string(tt.wantResponseBody) {
				t.Errorf("POST() gotResponseBody = %v, want %v", string(gotResponseBody), string(tt.wantResponseBody))
			} else {
				t.Logf("POST() gotResponseBody = %v, want %v", string(gotResponseBody), string(tt.wantResponseBody))
			}
		})
	}
}
func TestJsonFetch_PUT(t *testing.T) {
	shutdown := runEcho()
	defer shutdown()

	time.Sleep(time.Second)
	fetch := NewJsonFetch(nil, time.Second*5)

	tests := []struct {
		name             string
		fetch            *JsonFetch
		url              string
		wantResponseBody []byte
		wantErr          bool
	}{
		{"TestJsonFetch_PUT_001", fetch, _TEST_URL + "/get/200", []byte(`{"message":"ok"}`), false},
		{"TestJsonFetch_PUT_002", fetch, _TEST_URL + "/get/307", []byte(`{"message":"ok"}`), false},
		{"TestJsonFetch_PUT_003", fetch, _TEST_URL + "/get/400", []byte(`{"message":"failed"}`), true},
		{"TestJsonFetch_PUT_004", fetch, _TEST_URL + "/get/404", []byte(`{"message":"failed"}`), true},
		{"TestJsonFetch_PUT_005", fetch, _TEST_URL + "/get/500", []byte(`{"message":"failed"}`), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.fetch
			gotResponseBody, err := r.PUT(tt.url, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("PUT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(gotResponseBody) == string(tt.wantResponseBody) {
				t.Errorf("PUT() gotResponseBody = %v, want %v", string(gotResponseBody), string(tt.wantResponseBody))
			} else {
				t.Logf("PUT() gotResponseBody = %v, want %v", string(gotResponseBody), string(tt.wantResponseBody))
			}
		})
	}
}
func TestJsonFetch_PATCH(t *testing.T) {
	shutdown := runEcho()
	defer shutdown()

	time.Sleep(time.Second)
	fetch := NewJsonFetch(nil, time.Second*5)

	tests := []struct {
		name             string
		fetch            *JsonFetch
		url              string
		wantResponseBody []byte
		wantErr          bool
	}{
		{"TestJsonFetch_PATCH_001", fetch, _TEST_URL + "/get/200", []byte(`{"message":"ok"}`), false},
		{"TestJsonFetch_PATCH_002", fetch, _TEST_URL + "/get/307", []byte(`{"message":"ok"}`), false},
		{"TestJsonFetch_PATCH_003", fetch, _TEST_URL + "/get/400", []byte(`{"message":"failed"}`), true},
		{"TestJsonFetch_PATCH_004", fetch, _TEST_URL + "/get/404", []byte(`{"message":"failed"}`), true},
		{"TestJsonFetch_PATCH_005", fetch, _TEST_URL + "/get/500", []byte(`{"message":"failed"}`), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.fetch
			gotResponseBody, err := r.PATCH(tt.url, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("PATCH() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(gotResponseBody) == string(tt.wantResponseBody) {
				t.Errorf("PATCH() gotResponseBody = %v, want %v", string(gotResponseBody), string(tt.wantResponseBody))
			} else {
				t.Logf("PATCH() gotResponseBody = %v, want %v", string(gotResponseBody), string(tt.wantResponseBody))
			}
		})
	}
}
func TestJsonFetch_DELETE(t *testing.T) {
	shutdown := runEcho()
	defer shutdown()

	time.Sleep(time.Second)
	fetch := NewJsonFetch(nil, time.Second*5)

	tests := []struct {
		name             string
		fetch            *JsonFetch
		url              string
		wantResponseBody []byte
		wantErr          bool
	}{
		{"TestJsonFetch_DELETE_001", fetch, _TEST_URL + "/get/200", []byte(`{"message":"ok"}`), false},
		{"TestJsonFetch_DELETE_002", fetch, _TEST_URL + "/get/307", []byte(`{"message":"ok"}`), false},
		{"TestJsonFetch_DELETE_003", fetch, _TEST_URL + "/get/400", []byte(`{"message":"failed"}`), true},
		{"TestJsonFetch_DELETE_004", fetch, _TEST_URL + "/get/404", []byte(`{"message":"failed"}`), true},
		{"TestJsonFetch_DELETE_005", fetch, _TEST_URL + "/get/500", []byte(`{"message":"failed"}`), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.fetch
			gotResponseBody, err := r.DELETE(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("DELETE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(gotResponseBody) == string(tt.wantResponseBody) {
				t.Errorf("DELETE() gotResponseBody = %v, want %v", string(gotResponseBody), string(tt.wantResponseBody))
			} else {
				t.Logf("DELETE() gotResponseBody = %v, want %v", string(gotResponseBody), string(tt.wantResponseBody))
			}
		})
	}
}
func TestJsonFetch_CONNECT(t *testing.T) {
	shutdown := runEcho()
	defer shutdown()

	time.Sleep(time.Second)
	fetch := NewJsonFetch(nil, time.Second*5)

	tests := []struct {
		name             string
		fetch            *JsonFetch
		url              string
		wantResponseBody []byte
		wantErr          bool
	}{
		{"TestJsonFetch_CONNECT_001", fetch, _TEST_URL + "/get/200", []byte(`{"message":"ok"}`), false},
		{"TestJsonFetch_CONNECT_002", fetch, _TEST_URL + "/get/307", []byte(`{"message":"ok"}`), false},
		{"TestJsonFetch_CONNECT_003", fetch, _TEST_URL + "/get/400", []byte(`{"message":"failed"}`), true},
		{"TestJsonFetch_CONNECT_004", fetch, _TEST_URL + "/get/404", []byte(`{"message":"failed"}`), true},
		{"TestJsonFetch_CONNECT_005", fetch, _TEST_URL + "/get/500", []byte(`{"message":"failed"}`), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.fetch
			gotResponseBody, err := r.CONNECT(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("CONNECT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(gotResponseBody) == string(tt.wantResponseBody) {
				t.Errorf("CONNECT() gotResponseBody = %v, want %v", string(gotResponseBody), string(tt.wantResponseBody))
			} else {
				t.Logf("CONNECT() gotResponseBody = %v, want %v", string(gotResponseBody), string(tt.wantResponseBody))
			}
		})
	}
}
func TestJsonFetch_OPTIONS(t *testing.T) {
	shutdown := runEcho()
	defer shutdown()

	time.Sleep(time.Second)
	fetch := NewJsonFetch(nil, time.Second*5)

	tests := []struct {
		name             string
		fetch            *JsonFetch
		url              string
		wantResponseBody []byte
		wantErr          bool
	}{
		{"TestJsonFetch_OPTIONS_001", fetch, _TEST_URL + "/get/200", []byte(`{"message":"ok"}`), false},
		{"TestJsonFetch_OPTIONS_002", fetch, _TEST_URL + "/get/307", []byte(`{"message":"ok"}`), false},
		{"TestJsonFetch_OPTIONS_003", fetch, _TEST_URL + "/get/400", []byte(`{"message":"failed"}`), true},
		{"TestJsonFetch_OPTIONS_004", fetch, _TEST_URL + "/get/404", []byte(`{"message":"failed"}`), true},
		{"TestJsonFetch_OPTIONS_005", fetch, _TEST_URL + "/get/500", []byte(`{"message":"failed"}`), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.fetch
			gotResponseBody, err := r.OPTIONS(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("OPTIONS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(gotResponseBody) == string(tt.wantResponseBody) {
				t.Errorf("OPTIONS() gotResponseBody = %v, want %v", string(gotResponseBody), string(tt.wantResponseBody))
			} else {
				t.Logf("OPTIONS() gotResponseBody = %v, want %v", string(gotResponseBody), string(tt.wantResponseBody))
			}
		})
	}
}
func TestJsonFetch_TRACE(t *testing.T) {
	shutdown := runEcho()
	defer shutdown()

	time.Sleep(time.Second)
	fetch := NewJsonFetch(nil, time.Second*5)

	tests := []struct {
		name             string
		fetch            *JsonFetch
		url              string
		wantResponseBody []byte
		wantErr          bool
	}{
		{"TestJsonFetch_TRACE_001", fetch, _TEST_URL + "/get/200", []byte(`{"message":"ok"}`), false},
		{"TestJsonFetch_TRACE_002", fetch, _TEST_URL + "/get/307", []byte(`{"message":"ok"}`), false},
		{"TestJsonFetch_TRACE_003", fetch, _TEST_URL + "/get/400", []byte(`{"message":"failed"}`), true},
		{"TestJsonFetch_TRACE_004", fetch, _TEST_URL + "/get/404", []byte(`{"message":"failed"}`), true},
		{"TestJsonFetch_TRACE_005", fetch, _TEST_URL + "/get/500", []byte(`{"message":"failed"}`), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.fetch
			gotResponseBody, err := r.TRACE(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("TRACE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(gotResponseBody) == string(tt.wantResponseBody) {
				t.Errorf("TRACE() gotResponseBody = %v, want %v", string(gotResponseBody), string(tt.wantResponseBody))
			} else {
				t.Logf("TRACE() gotResponseBody = %v, want %v", string(gotResponseBody), string(tt.wantResponseBody))
			}
		})
	}
}
