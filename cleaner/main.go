package main

import (
	"encoding/json"
	"log"

	"github.com/guad/paperless2/backend/broker"
	"github.com/guad/paperless2/backend/model"
	"github.com/guad/paperless2/backend/storage"
	"github.com/streadway/amqp"
)

func main() {
	log.Println("Starting up..")

	storage.InitStorage()
	broker.InitBroker()

	s3q, err := broker.RabbitMQ.QueueDeclare(
		"s3_cleanup",
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
		s3q.Name,                    // queue name
		"",                          // routing key
		broker.DocumentCleanupQueue, // exchange
		false,                       // nowait
		nil,                         // args
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

	go s3CleanupWorker(msgs)

	inf := make(chan bool)

	<-inf

	log.Println("Shutting down..")
}

func s3CleanupWorker(queue <-chan amqp.Delivery) {
	for d := range queue {
		var doc model.Document

		_ = json.Unmarshal(d.Body, &doc)

		err := storage.S3.RemoveObject(storage.DocumentBucket, doc.S3Path)

		if err != nil {
			log.Println("CLEANUP ERROR:", err)
			d.Nack(false, true)
			continue
		}

		if doc.ThumbnailPath != "" {
			err = storage.S3.RemoveObject(storage.DocumentBucket, doc.ThumbnailPath)

			if err != nil {
				log.Println("THUMB CLEANUP ERROR:", err)
				// Don't NACK this error as it will requeue
				// and fail when removing the main document.
				// Just log it.
			}
		}

		log.Println("Cleaned up", doc.Filename)
		d.Ack(false)
	}
}
