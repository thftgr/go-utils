package handlefunc

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

func RequestHandler[T any](cb func(echo.Context, T) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req T
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return cb(c, req)
	}
}

func RequestResponseJsonHandler[T, R any](cb func(echo.Context, *T) (*R, error)) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req T
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if res, err := cb(c, &req); err != nil {
			return err
		} else {
			return c.JSON(http.StatusOK, res)
		}
	}
}

func RequestResponseXmlHandler[T, R any](cb func(echo.Context, *T) (*R, error)) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req T
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if res, err := cb(c, &req); err != nil {
			return err
		} else {
			return c.XML(http.StatusOK, res)
		}
	}
}

func RequestResponseHtmlHandler[T any](cb func(echo.Context, *T) (string, error)) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req T
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if res, err := cb(c, &req); err != nil {
			return err
		} else {
			return c.HTML(http.StatusOK, res)
		}
	}
}

func RequestResponseFilePathHandler[T any](cb func(echo.Context, *T) (string, error)) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req T
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if res, err := cb(c, &req); err != nil {
			return err
		} else {
			return c.File(res)
		}
	}
}
func RequestResponseFileReaderHandler[T any](cb func(echo.Context, *T) (reader io.ReadCloser, size int64, filename string, err error)) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req T
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if reader, size, filename, err := cb(c, &req); err != nil {
			return err

		} else {
			c.Response().Header().Set(echo.HeaderContentType, `application/octet-stream`)
			c.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf(`attachment; filename="%s"`, filename))
			c.Response().Header().Set(echo.HeaderContentLength, fmt.Sprintf("%d", size))
			defer func() {
				err = errors.Join(err, reader.Close())
			}()
			_, err = io.Copy(c.Response().Writer, reader)
			return err
		}
	}
}
func RequestResponseRedirectHandler[T any](cb func(echo.Context, *T) (string, error)) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req T
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if url, err := cb(c, &req); err != nil {
			return err
		} else {
			return c.Redirect(http.StatusTemporaryRedirect, url)
		}
	}
}
