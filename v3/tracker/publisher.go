package tracker

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (l *TrackerPublisher) SetAMQPConnection(connection *amqp.Connection) {
	l.connection = connection
}

func (t *TrackerPublisher) TrackUserSearch(tp *TrackerPayload) {
	err := t.pushToQueue(tp,SearchType)
	if err != nil {
		log.Printf("Track User Search Publisher Error: %s\n", err)
		return
	}
}

func (t *TrackerPublisher) TrackUserCW(tp *TrackerPayload) {
	err := t.pushToQueue(tp,CWType)
	if err != nil {
		log.Printf("Track User CW Publisher Error: %s\n", err)
		return
	}
}

func (t *TrackerPublisher) pushToQueue(tp *TrackerPayload,topic string) error {
	err := setup(t.connection)
	if err != nil {
		return err
	}

	j, _ := json.Marshal(tp)
	err = t.push(string(j),topic)
	if err != nil {
		return err
	}

	return nil
}

func  (t *TrackerPublisher) push(event string,topic string) error {
	channel, err := t.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	log.Println("Pushing to channel")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = channel.PublishWithContext(ctx,
		"tracker_topic",          // exchange
		topic, // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body: []byte(event),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

