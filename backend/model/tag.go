package model

import (
	"github.com/globalsign/mgo/bson"
)

type Tag struct {
	ID      bson.ObjectId `json:"id,omitempty" bson:"_id"`
	Name    string        `json:"name,omitempty" bson:"name"`
	Regex   string        `json:"regex,omitempty" bson:"regex"`
	Implies []string      `json:"implies,omitempty" bson:"implies"`
}
