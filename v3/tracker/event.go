package tracker

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
)

var CWType string = "tracker.user.cw"
var SearchType string = "tracker.user.search"

type TrackerPayload struct {
	Bid                 string  `json:"bid,omitempty"`
	Type                string  `json:"type,omitempty"`
	Phone               string  `json:"phone,omitempty"`
	Name                string  `json:"name,omitempty"`
	Email               string  `json:"email,omitempty"`
	UtmSource           string  `json:"utm_source,omitempty"`
	UtmMedium           string  `json:"utm_medium,omitempty"`
	UtmCampaignId       string  `json:"utm_campaign_id,omitempty"`
	UtmCampaignName     string  `json:"utm_campaign_name,omitempty"`
	Lat                 float64 `json:"lat,omitempty,omitempty"`
	Lng                 float64 `json:"lng,omitempty,omitempty"`
	FloorNumber         int32   `json:"floorNumber,omitempty"`
	PricePerSqFt        int32   `json:"pricePerSqFt,omitempty"`
	AreaInSqft          int32   `json:"areaInSqft,omitempty"`
	Address             string  `json:"address,omitempty"`
	IpAddress           string  `json:"ipAddress,omitempty"`
	City                string  `json:"city,omitempty"`
	CbSearchListSize    string  `json:"cb_search_list_size,omitempty"`
	GoSearchListSize    string  `json:"go_search_list_size,omitempty"`
	TotalSearchListSize string  `json:"total_search_list_size,omitempty"`
	ClicworthPrice      int32   `json:"clicworthPrice,omitempty"`
	LowClicworthPrice   int32   `json:"lowClicworthPrice,omitempty"`
	HighClicworthPrice  int32   `json:"highClicworthPrice,omitempty"`
	ConfidenceLevel     string  `json:"confidenceLevel,omitempty"`
	IsGoogleSearch      bool    `json:"isGoogleSearch,omitempty"`
	CBProjectId         string  `json:"cbprojectid,omitempty"`
}

type TrackerPublisher struct {
	connection *amqp.Connection
}

type TrackerConsumer struct {
	connection       *amqp.Connection
	client           *mongo.Client
	database         *mongo.Database
	searchCollection *mongo.Collection
	cwCollection     *mongo.Collection
}

type TrackerEntry struct {
	ID                  string    `bson:"_id,omitempty" json:"id,omitempty"`
	Bid                 string    `bson:"bid,omitempty" json:"bid,omitempty"`
	Type                string    `bson:"type,omitempty" json:"type,omitempty"`
	Phone               string    `bson:"phone,omitempty" json:"phone,omitempty"`
	Name                string    `bson:"name,omitempty" json:"name,omitempty"`
	Email               string    `bson:"email,omitempty" json:"email,omitempty"`
	UtmSource           string    `bson:"utm_source,omitempty" json:"utm_source,omitempty"`
	UtmMedium           string    `bson:"utm_medium,omitempty" json:"utm_medium,omitempty"`
	UtmCampaignId       string    `bson:"utm_campaign_id,omitempty" json:"utm_campaign_id,omitempty"`
	UtmCampaignName     string    `bson:"utm_campaign_name,omitempty" json:"utm_campaign_name,omitempty"`
	Lat                 float64   `bson:"lat,omitempty" json:"lat,omitempty"`
	Lng                 float64   `bson:"lng,omitempty" json:"lng,omitempty"`
	FloorNumber         int32     `bson:"floorNumber,omitempty" json:"floorNumber,omitempty"`
	PricePerSqFt        int32     `bson:"pricePerSqFt,omitempty" json:"pricePerSqFt,omitempty"`
	AreaInSqft          int32     `bson:"areaInSqft,omitempty" json:"areaInSqft,omitempty"`
	ClicworthPrice      int32     `bson:"clicworthPrice,omitempty" json:"clicworthPrice,omitempty"`
	LowClicworthPrice   int32     `bson:"lowClicworthPrice,omitempty" json:"lowClicworthPrice,omitempty"`
	HighClicworthPrice  int32     `bson:"highClicworthPrice,omitempty" json:"highClicworthPrice,omitempty"`
	ConfidenceLevel     string    `bson:"confidenceLevel,omitempty" json:"confidenceLevel,omitempty"`
	IsGoogleSearch      bool      `bson:"isGoogleSearch,omitempty" json:"isGoogleSearch,omitempty"`
	CBProjectId         string    `bson:"cbprojectid,omitempty" json:"cbprojectid,omitempty"`
	Address             string    `bson:"address,omitempty" json:"address,omitempty"`
	IpAddress           string    `bson:"ip_address,omitempty" json:"ip_address,omitempty"`
	City                string    `bson:"city,omitempty" json:"city,omitempty"`
	CbSearchListSize    string    `bson:"cb_search_list_size,omitempty" json:"cb_search_list_size,omitempty"`
	GoSearchListSize    string    `bson:"go_search_list_size,omitempty" json:"go_search_list_size,omitempty"`
	TotalSearchListSize string    `bson:"total_search_list_size,omitempty" json:"total_search_list_size,omitempty"`
	UpdatedAt           int64     `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	CreatedAt           time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"`
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
		"tracker_topic", //name
		"topic",         //type
		true,            //durable
		false,           //auto-delete
		false,           //internal
		false,           // no-wait
		nil,             //arguments
	)
}

func declareRandomeQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    //name
		false, //durable
		false, //delete when unsued ?
		true,  // exclisive
		false, //no-wait?
		nil,   //arguments
	)
}
