package user

import (
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"github.com/guad/paperless2/backend/crypto"
	"github.com/guad/paperless2/backend/model"

	"github.com/globalsign/mgo/bson"
	"github.com/guad/paperless2/backend/db"

	"github.com/labstack/echo"
)

var (
	ErrNotLoggedIn = errors.New("user was not logged in")
)

type LoginModel struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func Login(c echo.Context) error {
	var data LoginModel
	var user model.User

	err := c.Bind(&data)

	if err != nil {
		return err
	}

	sesh := db.Ctx()
	defer sesh.Close()

	col := sesh.DB("paperless").C("users")

	err = col.Find(bson.M{"email": data.Username}).One(&user)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, struct{}{})
	}

	password, _ := hex.DecodeString(user.PasswordHash)

	if err = crypto.ComparePasswords(data.Password, password); err != nil {
		return c.JSON(http.StatusUnauthorized, struct{}{})
	}

	token, err := GenerateAccessToken(user)

	if err != nil {
		return err
	}

	c.SetCookie(&http.Cookie{
		Name:     "session",
		Value:    token,
		Secure:   true,
		HttpOnly: true,

		Expires: time.Now().Add(24 * time.Hour),
	})

	return c.JSON(http.StatusOK, struct {
		Token string `json:"token,omitempty"`
	}{token})
}

func Register(c echo.Context) error {
	var data LoginModel

	err := c.Bind(&data)

	if err != nil {
		return err
	}

	sesh := db.Ctx()
	defer sesh.Close()

	col := sesh.DB("paperless").C("users")

	count, err := col.Find(bson.M{"email": data.Username}).Count()

	if err != nil {
		return err
	}

	if count > 0 {
		return c.JSON(http.StatusBadRequest, struct {
			Success bool   `json:"success,omitempty"`
			Reason  string `json:"reason,omitempty"`
		}{
			Success: false,
			Reason:  "email_taken",
		})
	}

	// TODO: Validate email and password

	user := model.User{
		Email:        data.Username,
		ID:           bson.NewObjectId(),
		RegisterDate: time.Now(),
		PasswordHash: hex.EncodeToString(crypto.HashPassword(data.Password)),
	}

	err = col.Insert(user)

	if err != nil {
		return err
	}

	token, err := GenerateAccessToken(user)

	if err != nil {
		return err
	}

	c.SetCookie(&http.Cookie{
		Name:     "session",
		Value:    token,
		Secure:   true,
		HttpOnly: true,

		Expires: time.Now().Add(24 * time.Hour),
	})

	return c.JSON(http.StatusOK, struct {
		Success bool   `json:"success,omitempty"`
		Token   string `json:"token,omitempty"`
	}{true, token})
}

func GetUserID(c echo.Context) string {
	userv := c.Request().Context().Value("userid")

	if userv == nil {
		panic(ErrNotLoggedIn)
	}

	if _, ok := userv.(string); !ok {
		panic(ErrNotLoggedIn)
	}

	return userv.(string)
}
