package api

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"github.com/guad/paperless2/backend/broker"
	"github.com/guad/paperless2/backend/db"
	"github.com/guad/paperless2/backend/model"
	"github.com/guad/paperless2/backend/storage"
	"github.com/minio/minio-go"
	"github.com/streadway/amqp"

	"github.com/globalsign/mgo/bson"

	"github.com/labstack/echo"

	"crypto/sha256"
)

func PushFile(c echo.Context) error {
	file, err := c.FormFile("document")
	id := bson.NewObjectId()

	if err != nil {
		return err
	}

	title := c.FormValue("title")

	content, err := file.Open()

	if err != nil {
		return err
	}

	defer content.Close()

	buffer := bytes.Buffer{}
	base64Buffer := bytes.Buffer{}
	hashreader, hashwriter := io.Pipe()

	tee := io.MultiWriter(&buffer, &base64Buffer, hashwriter)
	key := filepath.Join("documents", id.Hex(), file.Filename)

	hash := make(chan string, 1)

	go func() {
		io.Copy(tee, content)
		hashwriter.Close()
	}()

	go func() {
		h := sha256.New()

		io.Copy(h, hashreader)

		hash <- hex.EncodeToString(h.Sum(nil))
	}()

	// TODO: Hash should be of the encrypted file.
	hashhex := <-hash

	// Find duplicates
	sesh := db.Ctx()
	defer sesh.Close()

	col := sesh.DB("paperless").C("documents")

	count, err := col.Find(bson.M{
		"hash": hashhex,
	}).Count()

	if err != nil {
		return err
	}

	// This file already exists
	if count >= 1 {
		return c.JSON(http.StatusAccepted, struct{}{})
	}

	_, err = storage.S3.PutObject(
		storage.DocumentBucket,
		key,
		&buffer,
		file.Size,
		minio.PutObjectOptions{
			ContentType: file.Header.Get("Content-Type"),
		})

	if err != nil {
		return err
	}

	doc := model.Document{
		ID:          id,
		Title:       title,
		Filename:    file.Filename,
		ContentType: file.Header.Get("Content-Type"),
		Timestamp:   time.Now(),
		S3Path:      key,
		Hash:        hashhex,
	}

	err = col.Insert(doc)

	if err != nil {
		// TODO: Delete object if fail
		return err
	}

	b64 := base64.StdEncoding.EncodeToString(base64Buffer.Bytes())

	packet := struct {
		Document model.Document `json:"document,omitempty"`
		Data     string         `json:"data,omitempty"`
	}{
		Document: doc,
		Data:     b64,
	}

	jsonbytes, _ := json.Marshal(packet)

	broker.RabbitMQ.Publish(
		broker.DocumentUploadQueue, // Exchange
		"",                         // routing key
		false,                      // mandatory
		false,                      // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         jsonbytes,
			DeliveryMode: amqp.Persistent,
		},
	)

	return c.JSON(http.StatusCreated, doc)
}
