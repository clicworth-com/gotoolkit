package logs

import (
	"encoding/json"
	"fmt"
	"log"
	"log-service/data"
	"net/http"

	"github.com/clicworth-com/gotoolkit"
	amqp "github.com/rabbitmq/amqp091-go"
)

type LogConsumer struct {
	conn *amqp.Connection
	queueName string
}

func NewConsumer(conn *amqp.Connection) (LogConsumer, error) {
	consumer := LogConsumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return LogConsumer{}, err
	}

	return consumer, nil
}

func (consumer *LogConsumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

// type Payload struct {
// 	Name string `json:"name"`
// 	Data string `json:"data"`
// }

func (consumer *LogConsumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
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

			go handlePayload(payload)
		}
	}()

	fmt.Printf("Waiting for message [Exchange, Queue] [logs_topic, %s]\n", q.Name)
	<-forever

	return nil
}

func handlePayload(payload Payload)  {
	switch payload.Name {
	case "log","event":
		// log whatever we get
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	default:
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	}
	
}

func WriteLog(w http.ResponseWriter, r *http.Request) {
	var tools gotoolkit.Tools
	// read json into var
	var requestPayload JSONPayload
	_ = tools.ReadJSON(w, r, &requestPayload)

	// insert data
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		tools.ErrorJSON(w, err)
		return
	}

	resp := gotoolkit.JSONResponse{
		Error: false,
		Message: "logged",
	}

	tools.WriteJSON(w, http.StatusAccepted, resp)
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