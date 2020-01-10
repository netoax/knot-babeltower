package main

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type RegisterCommand struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

type Schema struct {
	SensorID  int    `json:"sensorId" validate:"gte=150"`
	ValueType int    `json:"valueType"`
	Unit      int    `json:"unit"`
	TypeID    int    `json:"typeId"`
	Name      string `json:"name"`
}

type UpdateSchemaCommand struct {
	ID     string   `json:"id"`
	Schema []Schema `json:"schema"`
}

type ListThingsCommand struct{}

func listenMessages() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs_topic", // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,              // queue name
		"device.registered", // routing key
		"FogOut",            // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
}

func sendData() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"FogIn", // name
		"topic", // type
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	// register := &RegisterCommand{
	// 	ID:   "18f8160184330541",
	// 	Name: "th12",
	// }

	// schema := UpdateSchemaCommand{
	// 	"18f8160184330541",
	// 	[]Schema{
	// 		{
	// 			SensorID:  0,
	// 			ValueType: 3,
	// 			Unit:      0,
	// 			TypeID:    65521,
	// 			Name:      "LED",
	// 		},
	// 	},
	// }

	ltc := ListThingsCommand{}
	pmsg, _ := json.Marshal(ltc)

	err = ch.Publish(
		"fogIn",       // exchange
		"device.list", // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			Headers: map[string]interface{}{
				"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NzYzMTcxNjEsImlhdCI6MTU3NjI4MTE2MSwiaXNzIjoibWFpbmZsdXgiLCJzdWIiOiJqb2FvQHRlc3QuY29tIn0.e_JNTTCMugoMeavu-EZPlVqcMdiyBWgdfaplgf17GDo",
			},
			ContentType: "text/plain",
			Body:        pmsg,
		})
	failOnError(err, "Failed to publish a message")
}

func main() {

	forever := make(chan bool)

	// listenMessages()
	sendData()

	<-forever
}
