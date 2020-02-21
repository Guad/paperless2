package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"log"
	"path/filepath"

	"github.com/globalsign/mgo/bson"
	"github.com/guad/paperless2/backend/broker"
	"github.com/guad/paperless2/backend/db"
	"github.com/guad/paperless2/backend/model"
	"github.com/guad/paperless2/backend/storage"
	"github.com/minio/minio-go"
	"github.com/streadway/amqp"
)

func setupThumbnailer() {
	s3q, err := broker.RabbitMQ.QueueDeclare(
		"document_thumbnail_attach",
		true,  // durable
		false, // delete when unsued
		false, // exclusive
		false, // nowait
		nil,   // args
	)

	if err != nil {
		panic(err)
	}

	err = broker.RabbitMQ.QueueBind(
		s3q.Name,                  // queue name
		"",                        // routing key
		DocumentThumbnailComplete, // exchange
		false,                     // nowait
		nil,                       // args
	)

	if err != nil {
		panic(err)
	}

	msgs, err := broker.RabbitMQ.Consume(
		s3q.Name,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	if err != nil {
		panic(err)
	}

	go documentThumbnailer(msgs)
}

func documentThumbnailer(queue <-chan amqp.Delivery) {
	for d := range queue {
		sesh := db.Ctx()
		col := sesh.DB("paperless").C("documents")

		var result struct {
			Document  model.Document `json:"document,omitempty"`
			Thumbnail string         `json:"thumbnail,omitempty"`
		}

		_ = json.Unmarshal(d.Body, &result)

		if n, err := col.FindId(result.Document.ID).Count(); err != nil || n == 0 {
			// Document has been deleted. Ignore.
			sesh.Close()
			d.Ack(false)
			log.Println("Document", result.Document.ID, "has been deleted. Ignoring message.")
			continue
		}

		thumbnailBytes, _ := base64.StdEncoding.DecodeString(result.Thumbnail)
		thumbnailBuf := bytes.NewBuffer(thumbnailBytes)

		path := filepath.Join("thumbnails", result.Document.ID.Hex(), "thumbnail.png")

		_, err := storage.S3.PutObject(
			storage.DocumentBucket,
			path,
			thumbnailBuf,
			int64(thumbnailBuf.Len()),
			minio.PutObjectOptions{
				ContentType: "image/png",
			},
		)

		if err != nil {
			log.Println("ERROR Uploading thumbnail:", err)
			d.Nack(false, true)
			sesh.Close()
			continue
		}

		err = col.UpdateId(result.Document.ID, bson.M{
			"$set": bson.M{
				"thumbnail_path": path,
			},
		})

		if err != nil {
			log.Println("ERROR Updating thumbnail:", err)
			d.Nack(false, true)
			sesh.Close()
			continue
		}

		sesh.Close()
		d.Ack(false)

	}
}
