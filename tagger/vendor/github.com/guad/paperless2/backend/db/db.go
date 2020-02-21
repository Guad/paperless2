package db

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/globalsign/mgo"
)

type DatabaseConfig struct {
	Address  string `json:"address,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

var (
	mongo *mgo.Session
)

func InitDB() {
	secretPath := "mongo.json"

	if altp, ok := os.LookupEnv("MONGO_SECRETS"); ok {
		secretPath = altp
	}

	f, err := ioutil.ReadFile(secretPath)

	if err != nil {
		panic(err)
	}

	var config DatabaseConfig

	err = json.Unmarshal(f, &config)

	if err != nil {
		panic(err)
	}

	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{config.Address},
		Username: config.Username,
		Password: config.Password,
		Timeout:  10 * time.Second,
		Database: "",
	})

	if err != nil {
		panic(err)
	}

	mongo = session

	createModel()
}

func Ctx() *mgo.Session {
	return mongo.Clone()
}

func createModel() {
	mongo.DB("paperless").C("documents").Create(&mgo.CollectionInfo{})
	mongo.DB("paperless").C("tags").Create(&mgo.CollectionInfo{})
	mongo.DB("paperless").C("documents").EnsureIndex(mgo.Index{
		Key: []string{"$text:title", "$text:content"},
	})

	mongo.DB("paperless").C("documents").EnsureIndex(mgo.Index{
		Key: []string{"hash:1"},
	})
}
