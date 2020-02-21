package main

import (
	"encoding/json"
	"log"

	"github.com/globalsign/mgo/bson"
	"github.com/guad/paperless2/backend/broker"
	"github.com/guad/paperless2/backend/db"
	"github.com/guad/paperless2/backend/model"
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
		var result struct {
			Document  model.Document `json:"document,omitempty"`
			Thumbnail string         `json:"thumbnail,omitempty"`
		}

		_ = json.Unmarshal(d.Body, &result)

		sesh := db.Ctx()
		col := sesh.DB("paperless").C("documents")

		err := col.UpdateId(result.Document.ID, bson.M{
			"$set": bson.M{
				"thumbnail_path": result.Thumbnail,
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
