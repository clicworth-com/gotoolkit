package tracker

import (
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
)

var CWType string = "tracker.user.cw"
var SearchType string = "tracker.user.search"

type TrackerPayload struct {
	Bid string `json:"bid"`
	Type string `json:"type"`
	Phone string `json:"phone"`
	Name string `json:"name"`
	Email string `json:"email"`
	UtmSource string `json:"utm_source"`
	UtmMedium string `json:"utm_medium"`
	UtmCampaignId string `json:"utm_campaign_id"`
	UtmCampaignName string `json:"utm_campaign_name"`
	Lat float64 `json:"lat,omitempty"`
	Lng float64 `json:"lng,omitempty"`
	FloorNumber int32 `json:"floorNumber,omitempty"`
	TotalPrice int32 `json:"totalPrice,omitempty"`
	PricePerSqFt int32 `json:"pricePerSqFt,omitempty"`
	AreaInSqft int32 `json:"areaInSqft,omitempty"`
	Address string `json:"address"`
	IpAddress string `json:"ipAddress"`
	City string `json:"city"`
	CbSearchListSize string `json:"cb_search_list_size"`
	GoSearchListSize string `json:"go_search_list_size"`
	TotalSearchListSize string `json:"total_search_list_size"`
}

type TrackerPublisher struct {
	connection *amqp.Connection
}

type TrackerConsumer struct {
	connection *amqp.Connection
	client *mongo.Client
	database *mongo.Database
	searchCollection *mongo.Collection
	cwCollection *mongo.Collection
}

type TrackerEntry struct {
	ID string `bson:"_id,omitempty" json:"id,omitempty"`
	Bid string `bson:"bid" json:"bid"`
	Type string `bson:"type" json:"type"`
	Phone string `bson:"phone" json:"phone"`
	Name string `bson:"name" json:"name"`
	Email string `bson:"email" json:"email"`
	UtmSource string `bson:"utm_source" json:"utm_source"`
	UtmMedium string `bson:"utm_medium" json:"utm_medium"`
	UtmCampaignId string `bson:"utm_campaign_id" json:"utm_campaign_id"`
	UtmCampaignName string `bson:"utm_campaign_name" json:"utm_campaign_name"`
	Lat float64 `bson:"lat,omitempty" json:"lat,omitempty"`
	Lng float64 `bson:"lng,omitempty" json:"lng,omitempty"`
	FloorNumber int32 `bson:"floorNumber,omitempty" json:"floorNumber,omitempty"`
	TotalPrice int32 `bson:"totalPrice,omitempty" json:"totalPrice,omitempty"`
	PricePerSqFt int32 `bson:"pricePerSqFt,omitempty" json:"pricePerSqFt,omitempty"`
	AreaInSqft int32 `bson:"areaInSqft,omitempty" json:"areaInSqft,omitempty"`
	Address string `bson:"address" json:"address"`
	IpAddress string `bson:"ip_address" json:"ip_address"`
	City string `bson:"city" json:"city"`
	CbSearchListSize string `bson:"cb_search_list_size" json:"cb_search_list_size"`
	GoSearchListSize string `bson:"go_search_list_size" json:"go_search_list_size"`
	TotalSearchListSize string `bson:"total_search_list_size" json:"total_search_list_size"`
	UpdatedAt int64 `bson:"updated_at" json:"updated_at"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

func setup(conn *amqp.Connection) error {
	log.Println("inside setup")
	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()
	return declareExchange(channel)
}

func declareExchange(ch *amqp.Channel) error {
	log.Println("inside declareExchange")
	return ch.ExchangeDeclare(
		"tracker_topic", //name
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