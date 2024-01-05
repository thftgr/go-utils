package middleware

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/thftgr/go-utils/http/header"
	"net/http"
)

func ContentTypeFilter(types ...header.MimeType) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ct := header.ParseContentType(c.Request().Header.Get(echo.HeaderContentType))
			if ct == nil {
				return echo.NewHTTPError(http.StatusBadRequest, "invalid media type")
			}
			for i := range types {
				if types[i].Equals(ct.MimeType) {
					return next(c)
				}
			}
			return echo.NewHTTPError(http.StatusUnsupportedMediaType, fmt.Errorf("unsupported media type %s", ct.MimeType))
		}
	}
}
