package logs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewConsumer(conn *amqp.Connection, mongo *mongo.Client) (LogConsumer, error) {
	consumer := LogConsumer{
		connection: conn,
		client: mongo,
		database: mongo.Database("logs"),
		collection: mongo.Database("logs").Collection("logs"),
	}

	err := setup(consumer.connection)
	if err != nil {
		return LogConsumer{}, err
	}

	return consumer, nil
}

func (consumer *LogConsumer) Listen() error {
	topics := []string{"log.INFO", "log.WARNING", "log.ERROR"}
	ch, err := consumer.connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareRandomeQueue(ch)
	if err != nil {
		return err
	}

	for _, s := range topics {
		ch.QueueBind(
			q.Name,
			s,
			"logs_topic",
			false,
			nil,
		)

		if err != nil {
			return err
		}

	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func ()  {
		for d:= range messages {
			var payload LogPayload
			_ = json.Unmarshal(d.Body, &payload)

			go consumer.handlePayload(payload)
		}
	}()

	fmt.Printf("Waiting for message [Exchange, Queue] [logs_topic, %s]\n", q.Name)
	<-forever

	return nil
}

func (consumer *LogConsumer) handlePayload(payload LogPayload)  {
	// insert data
	logEntry := LogEntry{
		Service: payload.Service,
		Data: payload.Data,
		LogLevel: payload.LogLevel,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := consumer.insert(logEntry)
	if err != nil {
		return
	}
}

func (consumer *LogConsumer) insert(entry LogEntry) error {
	_, err := consumer.collection.InsertOne(context.TODO(), LogEntry{
		Service: entry.Service,
		Data: entry.Data,
		LogLevel: entry.LogLevel,
		CreatedAt: entry.CreatedAt,
		UpdatedAt: entry.UpdatedAt,
	})
	if err != nil {
		log.Println("Error inserting into logs: ",err)
		return err
	}

	return nil
}