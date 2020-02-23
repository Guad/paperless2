package user

import (
	"time"

	"github.com/gbrlsnchs/jwt/v3"
	"github.com/guad/paperless2/backend/crypto"
	"github.com/guad/paperless2/backend/model"
)

func GenerateAccessToken(user model.User) (string, error) {
	now := time.Now()
	rs256 := jwt.NewHS256(crypto.Key())

	p := jwt.Payload{
		Issuer:         "paperless2.kolhos.chichasov.es",
		Subject:        user.ID.Hex(),
		Audience:       jwt.Audience{"paperless2.kolhos.chichasov.es"},
		ExpirationTime: jwt.NumericDate(now.Add(24 * time.Hour)),
		IssuedAt:       jwt.NumericDate(now),
	}

	token, err := jwt.Sign(p, rs256)

	if err != nil {
		panic(err)
	}

	return string(token), nil
}

func ValidateToken(token string) (string, error) {
	now := time.Now()
	rs256 := jwt.NewHS256(crypto.Key())

	// First of all, decode code JWT.
	var p jwt.Payload

	iatValidator := jwt.IssuedAtValidator(now)
	expValidator := jwt.ExpirationTimeValidator(now)
	audValidator := jwt.AudienceValidator(jwt.Audience{"paperless2.kolhos.chichasov.es"})

	validatePayload := jwt.ValidatePayload(&p, iatValidator, expValidator, audValidator)

	_, err := jwt.Verify([]byte(token), rs256, &p, validatePayload)

	if err != nil {
		return "", err
	}

	return p.Subject, nil
}
