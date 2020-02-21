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

const (
	DocumentOCRComplete       = "document_ocr_complete"
	DocumentTaggingComplete   = "document_tagging_complete"
	DocumentThumbnailComplete = "document_thumbnail_complete"
)

func main() {
	log.Println("Starting up..")

	broker.InitBroker()
	db.InitDB()

	setupThumbnailer()

	err := broker.RabbitMQ.ExchangeDeclare(
		DocumentTaggingComplete,
		"fanout",
		true,  // durable
		false, // autodelete
		false, // internal
		false, // nowait
		nil,   // args
	)

	if err != nil {
		panic(err)
	}

	s3q, err := broker.RabbitMQ.QueueDeclare(
		"document_tagger",
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
		s3q.Name,            // queue name
		"",                  // routing key
		DocumentOCRComplete, // exchange
		false,               // nowait
		nil,                 // args
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

	go documentTagWorker(msgs)

	inf := make(chan bool)

	<-inf

	log.Println("Shutting down..")
}

func documentTagWorker(queue <-chan amqp.Delivery) {
	for d := range queue {
		var result struct {
			Document model.Document `json:"document,omitempty"`
			Content  string         `json:"content,omitempty"`
		}

		_ = json.Unmarshal(d.Body, &result)

		// Fetch all tags
		sesh := db.Ctx()

		col := sesh.DB("paperless").C("tags")

		var tags []model.Tag

		err := col.Find(bson.M{}).All(&tags)

		if err != nil {
			d.Nack(false, true)
			log.Println("Can't fetch tags:", err)
			sesh.Close()
			continue
		}

		doctags := tag(result.Content, tags)
		doccol := sesh.DB("paperless").C("documents")

		err = doccol.UpdateId(result.Document.ID, bson.M{
			"$set": bson.M{
				"tags":    doctags,
				"content": result.Content,
			},
		})

		if err != nil {
			d.Nack(false, true)
			log.Println("Can't fetch tags:", err)
			sesh.Close()
			continue
		}

		log.Println("Successfuly tagged", result.Document.ID)

		sesh.Close()

		d.Ack(false)
	}
}
