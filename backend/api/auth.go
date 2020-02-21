package api

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gbrlsnchs/jwt/v3"
	"github.com/labstack/echo"
)

type LoginModel struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func Login(c echo.Context) error {
	var model LoginModel
	err := c.Bind(&model)

	if err != nil {
		return err
	}

	// TODO: Improve

	usr := os.Getenv("ADMIN_USERNAME")
	pass := os.Getenv("ADMIN_PASSWORD")

	if pass != model.Password || usr != model.Username {
		return c.JSON(http.StatusUnauthorized, struct{}{})
	}

	token, err := GenerateAccessToken()

	if err != nil {
		return err
	}

	// TODO: Secure this
	c.SetCookie(&http.Cookie{
		Name:    "session",
		Value:   token,
		Expires: time.Now().Add(15 * 24 * time.Hour),
	})

	return c.JSON(http.StatusOK, struct {
		Token string `json:"token,omitempty"`
	}{token})
}

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

		if err := ValidateToken(token); err != nil {
			return c.String(http.StatusUnauthorized, "{}")
		}

		return next(c)
	}
}

func GenerateAccessToken() (string, error) {
	now := time.Now()
	rs256 := jwt.NewHS256([]byte(os.Getenv("ADMIN_KEY")))

	p := jwt.Payload{
		Issuer:         "paperless2.kolhos.chichasov.es",
		Subject:        "usr",
		Audience:       jwt.Audience{"paperless2.kolhos.chichasov.es"},
		ExpirationTime: jwt.NumericDate(now.Add(365 * 24 * time.Hour)),
		IssuedAt:       jwt.NumericDate(now),
	}

	token, err := jwt.Sign(p, rs256)

	if err != nil {
		panic(err)
	}

	return string(token), nil
}

func ValidateToken(token string) error {
	now := time.Now()
	rs256 := jwt.NewHS256([]byte(os.Getenv("ADMIN_KEY")))

	// First of all, decode code JWT.
	var p jwt.Payload

	iatValidator := jwt.IssuedAtValidator(now)
	expValidator := jwt.ExpirationTimeValidator(now)
	audValidator := jwt.AudienceValidator(jwt.Audience{"paperless2.kolhos.chichasov.es"})

	validatePayload := jwt.ValidatePayload(&p, iatValidator, expValidator, audValidator)

	_, err := jwt.Verify([]byte(token), rs256, &p, validatePayload)

	if err != nil {
		return err
	}

	return nil
}
