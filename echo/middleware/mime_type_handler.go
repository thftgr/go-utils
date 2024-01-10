package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/thftgr/go-utils/http/header"
	"net/http"
)

/*
   다음은 HTTP 메소드별로 `Accept` 및 `Content-Type` 헤더의 필요성과 그 처리 방법에 대한 비교 표입니다:

   | HTTP 메소드 || Accept 헤더 필요성    | 처리 방식 (Accept 헤더 누락 시)     || Content-Type 헤더 필요성   | 처리 방식 (Content-Type 헤더 누락 시)    |
   |-------------||--------------------|---------------------------------||--------------------------|---------------------------------------|
   | GET         || 선택적              | `* / *` 처리 또는 기본값 사용       || 일반적으로 불필요           | 오류 반환 없음 (본문 없음)                |
   | HEAD        || 선택적              | `* / *` 처리 또는 기본값 사용       || 일반적으로 불필요           | 오류 반환 없음 (본문 없음)                |
   | POST        || 선택적              | `* / *` 처리 또는 기본값 사용       || 필요 (본문이 있을 경우)     | 오류 반환 또는 기본값 사용 가능             |
   | PUT         || 선택적              | `* / *` 처리 또는 기본값 사용       || 필요 (본문이 있을 경우)     | 오류 반환 또는 기본값 사용 가능             |
   | DELETE      || 선택적              | `* / *` 처리 또는 기본값 사용       || 일반적으로 불필요           | 오류 반환 없음 (대부분 본문 없음)          |
   | PATCH       || 선택적              | `* / *` 처리 또는 기본값 사용       || 필요 (본문이 있을 경우)     | 오류 반환 또는 기본값 사용 가능             |
   | OPTIONS     || 선택적              | `* / *` 처리 또는 기본값 사용       || 일반적으로 불필요           | 오류 반환 없음 (대부분 본문 없음)          |
   |-------------||--------------------|---------------------------------||--------------------------|---------------------------------------|

   - Accept 헤더:
       대부분의 HTTP 메소드에서 `Accept` 헤더는 선택적입니다.
       헤더가 누락된 경우 일반적으로 클라이언트가 모든 MIME 타입을 수용한다고 간주하고 `* / *`로 처리하거나 서버 설정에 따른 기본값을 사용합니다.
   - Content-Type 헤더:
       `POST`, `PUT`, `PATCH`와 같이 요청 본문을 포함하는 메소드에서는 `Content-Type` 헤더가 필요합니다.
       이 헤더가 누락된 경우, 서버는 오류를 반환하거나 설정에 따라 기본 MIME 타입을 사용할 수 있습니다.
       그러나 `GET`, `HEAD`, `DELETE`, `OPTIONS`와 같이 본문이 없는 요청에서는 이 헤더가 필요하지 않습니다.
*/

// CheckAllowedAccept 클라이언트가 받을수 있는 타입에 서버가 처리 가능한 mimeTypes 가 포함되는지 검사
func CheckAllowedAccept(mimeTypes ...header.MimeType) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var at []header.Accept
			var mts []header.MimeType

			ha := c.Request().Header.Get(echo.HeaderAccept)
			if ha == "" {
				mts = append(mts, "*/*")
			} else {
				at = header.ParseAccept(ha)
				for i := range at {
					if at[i].MimeType.Valid() {
						mts = append(mts, at[i].MimeType)
					}
				}
			}

			for i := range mts {
				if mt := mts[i].Match(mimeTypes...); mt != nil {
					return next(c)
				}
			}
			return echo.NewHTTPError(http.StatusNotAcceptable) // 406 Not Acceptable
		}
	}
}

// CheckAllowedContentType 서버가 받을수 있는 타입에 client 의 MimeType 이 포함되는지 검사
func CheckAllowedContentType(mimeTypes ...header.MimeType) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ct := header.ParseContentType(c.Request().Header.Get(echo.HeaderContentType))
			if ct != nil {
				for i := range mimeTypes {
					if mt := mimeTypes[i].Match(ct.MimeType); mt != nil {
						return next(c)
					}
				}
			}

			return echo.NewHTTPError(http.StatusUnsupportedMediaType) // 415 Unsupported Media Type
		}
	}
}
