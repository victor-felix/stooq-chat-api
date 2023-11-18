package models

type Config struct {
	DevMode bool `split_words:"true" required:"true" default:"true"`
	MongoURL string `split_words:"true" required:"true" default:"mongodb://guest:guest@localhost:27017"`
	DatabaseName string `split_words:"true" required:"true" default:"chat"`
	JwtSecret string `split_words:"true" required:"true" default:"f5d9e7a587e6efbbbb8efbe71e6dd1f42cd6f040"`
	JwtTTL int `split_words:"true" required:"true" default:"15000"`
	Server struct {
		Port string `split_words:"true" required:"true" default:"5001"`
	}
	Websocket struct {
		ReadBufferSize int `split_words:"true" required:"true" default:"1024"`
		WriteBufferSize int `split_words:"true" required:"true" default:"1024"`
	}
	RabbitMQ struct {
		DSN string `split_words:"true" required:"true" default:"amqp://guest:guest@localhost:5672"`
		StooqReceiverQueue string `split_words:"true" required:"true" default:"stock-result"`
		StooqPublisherQueue string `split_words:"true" required:"true" default:"stock-process-request"`
	}
	StooqBaseUrl string `split_words:"true" required:"true" default:"https://stooq.com"`
}
