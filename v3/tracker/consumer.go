package tracker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewConsumer(conn *amqp.Connection, mongo *mongo.Client) (TrackerConsumer, error) {
	consumer := TrackerConsumer{
		connection:       conn,
		client:           mongo,
		database:         mongo.Database("tracker"),
		searchCollection: mongo.Database("tracker").Collection("search"),
		cwCollection:     mongo.Database("tracker").Collection("clicworth"),
		txCollection:     mongo.Database("tracker").Collection("txerrors"),
	}

	err := setup(consumer.connection)
	if err != nil {
		return TrackerConsumer{}, err
	}

	return consumer, nil
}

func (consumer *TrackerConsumer) Listen() error {
	topics := []string{CWType, SearchType, TxType}
	ch, err := consumer.connection.Channel()
	if err != nil {
		log.Printf("TrackerConsumer channel error %s", err)
		return err
	}
	defer ch.Close()

	q, err := declareRandomeQueue(ch)
	if err != nil {
		log.Printf("TrackerConsumer declareRandomeQueue error %s", err)
		return err
	}

	for _, s := range topics {
		ch.QueueBind(
			q.Name,
			s,
			"tracker_topic",
			false,
			nil,
		)
		if err != nil {
			log.Printf("TrackerConsumer QueueBind error %s", err)
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Printf("TrackerConsumer Consume	 error %s", err)
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			if d.RoutingKey == TxType {
				var payload TxTracker
				_ = json.Unmarshal(d.Body, &payload)
				go consumer.handleTxPayload(payload)
			} else {
				var payload TrackerPayload
				_ = json.Unmarshal(d.Body, &payload)
				go consumer.handlePayload(payload)
			}
		}
	}()

	fmt.Printf("Waiting for message [Exchange, Queue] [tracker_topic, %s]\n", q.Name)
	<-forever

	return nil
}

func (consumer *TrackerConsumer) handleTxPayload(payload TxTracker) {
	_, err := consumer.txCollection.InsertOne(context.TODO(), payload)
	if err != nil {
		log.Println("Error inserting into insertTxEntry: ", err)
	}
}

func (consumer *TrackerConsumer) handlePayload(payload TrackerPayload) {
	// insert data
	if payload.Type == SearchType {
		trackerEntry := TrackerEntry{
			Bid:                 payload.Bid,
			Type:                payload.Type,
			Phone:               payload.Phone,
			Name:                payload.Name,
			Email:               payload.Email,
			UtmSource:           payload.UtmSource,
			UtmMedium:           payload.UtmMedium,
			UtmCampaignId:       payload.UtmCampaignId,
			UtmCampaignName:     payload.UtmCampaignName,
			Address:             payload.Address,
			IpAddress:           payload.IpAddress,
			City:                payload.City,
			CbSearchListSize:    payload.CbSearchListSize,
			GoSearchListSize:    payload.GoSearchListSize,
			TotalSearchListSize: payload.TotalSearchListSize,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now().Unix(),
		}
		err := consumer.insertSearchEntry(trackerEntry)
		if err != nil {
			log.Printf("TrackerConsumer insertSearchEntry error %s", err)
			return
		}
	} else if payload.Type == CWType {
		trackerEntry := TrackerEntry{
			Bid:                 payload.Bid,
			Type:                payload.Type,
			Phone:               payload.Phone,
			Name:                payload.Name,
			Email:               payload.Email,
			UtmSource:           payload.UtmSource,
			UtmMedium:           payload.UtmMedium,
			UtmCampaignId:       payload.UtmCampaignId,
			UtmCampaignName:     payload.UtmCampaignName,
			Lat:                 payload.Lat,
			Lng:                 payload.Lng,
			FloorNumber:         payload.FloorNumber,
			ClicworthPrice:      payload.ClicworthPrice,
			LowClicworthPrice:   payload.LowClicworthPrice,
			HighClicworthPrice:  payload.HighClicworthPrice,
			ConfidenceLevel:     payload.ConfidenceLevel,
			IsGoogleSearch:      payload.IsGoogleSearch,
			CBProjectId:         payload.CBProjectId,
			PricePerSqFt:        payload.PricePerSqFt,
			AreaInSqft:          payload.AreaInSqft,
			Address:             payload.Address,
			IpAddress:           payload.IpAddress,
			City:                payload.City,
			CbSearchListSize:    payload.CbSearchListSize,
			GoSearchListSize:    payload.GoSearchListSize,
			TotalSearchListSize: payload.TotalSearchListSize,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now().Unix(),
		}
		err := consumer.insertCWEntry(trackerEntry)
		if err != nil {
			log.Printf("TrackerConsumer insertCWEntry error %s", err)
			return
		}
	}
}

func (consumer *TrackerConsumer) insertCWEntry(payload TrackerEntry) error {
	_, err := consumer.cwCollection.InsertOne(context.TODO(), TrackerEntry{
		Bid:                 payload.Bid,
		Type:                payload.Type,
		Phone:               payload.Phone,
		Name:                payload.Name,
		Email:               payload.Email,
		UtmSource:           payload.UtmSource,
		UtmMedium:           payload.UtmMedium,
		UtmCampaignId:       payload.UtmCampaignId,
		UtmCampaignName:     payload.UtmCampaignName,
		Lat:                 payload.Lat,
		Lng:                 payload.Lng,
		FloorNumber:         payload.FloorNumber,
		ClicworthPrice:      payload.ClicworthPrice,
		LowClicworthPrice:   payload.LowClicworthPrice,
		HighClicworthPrice:  payload.HighClicworthPrice,
		ConfidenceLevel:     payload.ConfidenceLevel,
		IsGoogleSearch:      payload.IsGoogleSearch,
		CBProjectId:         payload.CBProjectId,
		PricePerSqFt:        payload.PricePerSqFt,
		AreaInSqft:          payload.AreaInSqft,
		Address:             payload.Address,
		IpAddress:           payload.IpAddress,
		City:                payload.City,
		CbSearchListSize:    payload.CbSearchListSize,
		GoSearchListSize:    payload.GoSearchListSize,
		TotalSearchListSize: payload.TotalSearchListSize,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now().Unix(),
	})
	if err != nil {
		log.Println("Error inserting into insertCWEntry: ", err)
		return err
	}

	return nil
}

func (consumer *TrackerConsumer) insertSearchEntry(entry TrackerEntry) error {
	filter := bson.D{{Key: "bid", Value: entry.Bid}}
	filter = append(filter, bson.E{Key: "updated_at", Value: bson.D{
		{Key: "$gte", Value: time.Now().Unix() - 20},
	}})

	var result TrackerEntry
	err := consumer.searchCollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			_, err = consumer.searchCollection.InsertOne(context.TODO(), TrackerEntry{
				Bid:                 entry.Bid,
				Type:                entry.Type,
				Phone:               entry.Phone,
				Name:                entry.Name,
				Email:               entry.Email,
				UtmSource:           entry.UtmSource,
				UtmMedium:           entry.UtmMedium,
				UtmCampaignId:       entry.UtmCampaignId,
				UtmCampaignName:     entry.UtmCampaignName,
				Address:             entry.Address,
				IpAddress:           entry.IpAddress,
				City:                entry.City,
				CbSearchListSize:    entry.CbSearchListSize,
				GoSearchListSize:    entry.GoSearchListSize,
				TotalSearchListSize: entry.TotalSearchListSize,
				CreatedAt:           time.Now(),
				UpdatedAt:           time.Now().Unix(),
			})
			if err != nil {
				log.Println("Error inserting user search entry: ", err)
				return err
			}
		}
	} else {
		newFilter := bson.D{{Key: "_id", Value: result.ID}}
		update := bson.D{{Key: "$set", Value: bson.D{
			{Key: "address", Value: entry.Address},
			{Key: "cb_search_list_size", Value: entry.CbSearchListSize},
			{Key: "go_search_list_size", Value: entry.GoSearchListSize},
			{Key: "total_search_list_size", Value: entry.TotalSearchListSize},
			{Key: "updated_at", Value: time.Now().Unix()},
		}}}
		_, err := consumer.searchCollection.UpdateOne(context.TODO(), newFilter, update)
		if err != nil {
			log.Println("Error Updating user search entry: ", err)
			return err
		}
	}
	return nil
}
