package logs

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (l *LogPublisher) SetLogLevel(level uint) {
	switch level {
	case 0:
		l.logInfoLevelEnabled = false
		l.logWarningLevelEnabled = false
		l.logErrorLevelEnabled = true
	case 1:
		l.logInfoLevelEnabled = false
		l.logWarningLevelEnabled = true
		l.logErrorLevelEnabled = true
	case 2:
		l.logInfoLevelEnabled = true
		l.logWarningLevelEnabled = true
		l.logErrorLevelEnabled = true
	default:
		l.logInfoLevelEnabled = false
		l.logWarningLevelEnabled = false
		l.logErrorLevelEnabled = true
	}
}

func (l *LogPublisher) SetAMQPConnection(connection *amqp.Connection) {
	l.connection = connection
}

func (l *LogPublisher) LogINFO(serviceName,data string) {
	if l.logInfoLevelEnabled {
		err := l.pushToQueue(serviceName, data, "log.INFO")
		if err != nil {
			log.Printf("LogINFO %s  Error: %s\n", serviceName, err)
			return
		}
	}
}

func (l *LogPublisher) LogWARNING(serviceName,data string) {
	if l.logWarningLevelEnabled {
		err := l.pushToQueue(serviceName, data, "log.WARNING")
		if err != nil {
			log.Printf("LogWARNING %s  Error: %s\n", serviceName, err)
			return
		}
	}
}
func (l *LogPublisher) LogERROR(serviceName,data string) {
	if l.logWarningLevelEnabled {
		err := l.pushToQueue(serviceName, data, "log.ERROR")
		if err != nil {
			log.Printf("LogERROR %s  Error: %s\n", serviceName, err)
			return
		}
	}
}

func (l *LogPublisher) pushToQueue(serviceName, msg, severity string) error {
	err := setup(l.connection)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Service:serviceName,
		Data: msg,
		LogLevel: severity,
	}

	j, _ := json.Marshal(&payload)
	err = l.push(string(j), severity)
	if err != nil {
		return err
	}

	return nil
}

func (l *LogPublisher) push(event string, severity string) error {
	channel, err := l.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	log.Println("Pushing to channel")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = channel.PublishWithContext(ctx,
		"logs_topic",          // exchange
		severity, // routing key
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

