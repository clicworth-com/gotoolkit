package tracker

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (t *TrackerPublisher) SetAMQPConnection(connection *amqp.Connection) {
	t.connection = connection
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
		log.Printf("TrackerPublisher pushToQueue setup error %s", err)
		return err
	}

	j, _ := json.Marshal(tp)
	err = t.push(string(j),topic)
	if err != nil {
		log.Printf("TrackerPublisher pushToQueue push error %s", err)
		return err
	}

	return nil
}

func  (t *TrackerPublisher) push(event string,topic string) error {
	channel, err := t.connection.Channel()
	if err != nil {
		log.Printf("TrackerPublisher push Channel error %s", err)
		return err
	}
	defer channel.Close()

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
		log.Printf("TrackerPublisher push PublishWithContext error %s", err)
		return err
	}

	return nil
}

