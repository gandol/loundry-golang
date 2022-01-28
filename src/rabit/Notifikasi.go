package rabit

import (
	"github.com/streadway/amqp"
	"log"
	"loundry/api/src/helper"
)

func ReadMessage() {
	amqpServer := helper.ReadEnv("AMQP_SERVER")
	connectRabitMq, err := amqp.Dial(amqpServer)
	if err != nil {
		panic(err)
	}
	defer connectRabitMq.Close()
	amqpChannel, err := connectRabitMq.Channel()
	if err != nil {
		panic(err)
	}
	defer amqpChannel.Close()

	messages, err := amqpChannel.Consume(
		"notification", // queue
		"",             // consumer
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	if err != nil {
		panic(err)
	}
	forever := make(chan bool)
	go func() {
		for message := range messages {
			log.Printf("Received a message: %s", message.Body)
		}
	}()
	<-forever
}
