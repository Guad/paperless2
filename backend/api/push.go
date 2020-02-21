package api

import (
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

	s3reader, s3writer := io.Pipe()
	hashreader, hashwriter := io.Pipe()

	tee := io.MultiWriter(s3writer, hashwriter)
	key := filepath.Join("documents", id.Hex(), file.Filename)

	hash := make(chan string, 1)

	go func() {
		io.Copy(tee, content)
		s3writer.Close()
		hashwriter.Close()
	}()

	go func() {
		h := sha256.New()

		io.Copy(h, hashreader)

		hash <- hex.EncodeToString(h.Sum(nil))
	}()

	_, err = storage.S3.PutObject(
		storage.DocumentBucket,
		key,
		s3reader,
		file.Size,
		minio.PutObjectOptions{
			ContentType: file.Header.Get("Content-Type"),
		})

	if err != nil {
		return err
	}

	sesh := db.Ctx()
	defer sesh.Close()

	col := sesh.DB("paperless").C("documents")

	doc := model.Document{
		ID:          id,
		Title:       title,
		Filename:    file.Filename,
		ContentType: file.Header.Get("Content-Type"),
		Timestamp:   time.Now(),
		S3Path:      key,
		Hash:        <-hash,
	}

	err = col.Insert(doc)

	if err != nil {
		// TODO: Delete object if fail
		return err
	}

	// TODO: Put on broker

	jsonbytes, _ := json.Marshal(doc)

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
