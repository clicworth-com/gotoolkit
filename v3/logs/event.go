package logs

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
)

type LogPayload struct {
	Service string `json:"service"`
	Data string `json:"data"`
	LogLevel string `json:"log_level"`
}

type LogPublisher struct {
	connection *amqp.Connection
	logInfoLevelEnabled bool
	logWarningLevelEnabled bool
	logErrorLevelEnabled bool
}

type LogConsumer struct {
	connection *amqp.Connection
	client *mongo.Client
	database *mongo.Database
	collection *mongo.Collection
}

type LogEntry struct {
	ID string `bson:"_id,omitempty" json:"id,omitempty"`
	Service string `bson:"service" json:"service"`
	Data string `bson:"data" json:"data"`
	LogLevel string `bson:"log_level" json:"log_level"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func setup(conn *amqp.Connection) error {
	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()
	return declareExchange(channel)
}

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", //name
		"topic", //type
		true, //durable
		false, //auto-delete
		false, //internal
		false, // no-wait
		nil, //arguments
	)
}

func declareRandomeQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"", //name
		false, //durable
		false, //delete when unsued ?
		true, // exclisive
		false, //no-wait?
		nil, //arguments
	)
}