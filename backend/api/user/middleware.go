package user

import (
	"context"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get("Authorization")
		cookie, err := c.Cookie("session")

		var token string

		if err == nil && cookie.Value != "" {
			token = cookie.Value
		} else {
			if header == "" || !strings.HasPrefix(header, "Bearer ") {
				return c.String(http.StatusUnauthorized, "{}")
			}

			token = header[7:]
		}

		if userid, err := ValidateToken(token); err == nil {
			ctx := c.Request().Context()
			newctx := context.WithValue(ctx, "userid", userid)

			c.SetRequest(c.Request().WithContext(newctx))

			return next(c)
		}

		return c.String(http.StatusUnauthorized, "{}")
	}
}
