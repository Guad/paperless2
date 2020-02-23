package model

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

type User struct {
	ID           bson.ObjectId `json:"id,omitempty" bson:"_id"`
	Email        string        `json:"email,omitempty" bson:"email"`
	PasswordHash string        `json:"-,omitempty" bson:"password_hash"`
	RegisterDate time.Time     `json:"register_date,omitempty" bson:"register_date"`
}
