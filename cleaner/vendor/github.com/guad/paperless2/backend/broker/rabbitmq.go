package broker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/streadway/amqp"
)

type BrokerConfig struct {
	Host     string `json:"host,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

var (
	RabbitMQ *amqp.Channel

	DocumentUploadQueue  = "document_created"
	DocumentCleanupQueue = "document_deleted"
)

func InitBroker() {
	secretPath := "/config/rabbitmq.json"

	if altp, ok := os.LookupEnv("RABBITMQ_SECRETS"); ok {
		secretPath = altp
	}

	f, err := ioutil.ReadFile(secretPath)

	if err != nil {
		panic(err)
	}

	var config BrokerConfig

	err = json.Unmarshal(f, &config)

	if err != nil {
		panic(err)
	}

	connString := fmt.Sprintf("amqp://%v:%v@%v/",
		config.Username,
		config.Password,
		config.Host,
	)

	conn, err := amqp.Dial(connString)

	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()

	if err != nil {
		panic(err)
	}

	RabbitMQ = ch

	declareQueues()

	errq := make(chan *amqp.Error, 1)
	errq = conn.NotifyClose(errq)

	go func() {
		err := <-errq
		log.Fatal("Connection to RabbitMQ severed:", err)
	}()
}

func declareQueues() {
	err := RabbitMQ.ExchangeDeclare(
		DocumentUploadQueue,
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

	err = RabbitMQ.ExchangeDeclare(
		DocumentCleanupQueue,
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
}
