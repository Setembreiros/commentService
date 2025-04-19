package provider

import (
	"commentservice/infrastructure/database/atlas"
	"commentservice/infrastructure/database/sql_db"
	"commentservice/infrastructure/kafka"
	"commentservice/internal/api"
	"commentservice/internal/bus"
	database "commentservice/internal/db"
	"commentservice/internal/feature/create_comment"
	"commentservice/internal/feature/delete_comment"
	"commentservice/internal/service"
)

type Provider struct {
	env     string
	connStr string
}

func NewProvider(env, connStr string) *Provider {
	return &Provider{
		env:     env,
		connStr: connStr,
	}
}

func (p *Provider) ProvideAtlasCLient() (*atlas.AtlasClient, error) {
	return atlas.NewAtlasClient(p.connStr)
}

func (p *Provider) ProvideDb() (*sql_db.SqlDatabase, error) {
	return sql_db.NewDatabase(p.connStr)
}

func (p *Provider) ProvideEventBus() (*bus.EventBus, error) {
	kafkaProducer, err := kafka.NewKafkaProducer(p.kafkaBrokers())
	if err != nil {
		return nil, err
	}

	return bus.NewEventBus(kafkaProducer), nil
}

func (p *Provider) ProvideApiEndpoint(sqlClient *sql_db.SqlDatabase, bus *bus.EventBus) *api.Api {
	return api.NewApiEndpoint(p.env, p.ProvideApiControllers(sqlClient, bus))
}

func (p *Provider) ProvideApiControllers(sqlClient *sql_db.SqlDatabase, bus *bus.EventBus) []api.Controller {
	return []api.Controller{
		create_comment.NewCreateCommentController(create_comment.NewCreateCommentService(service.GetTimeServiceInstance(), create_comment.NewCreateCommentRepository(database.NewDatabase(sqlClient)), bus)),
		delete_comment.NewDeleteCommentController(delete_comment.NewDeleteCommentService(delete_comment.NewDeleteCommentRepository(database.NewDatabase(sqlClient)), bus)),
	}
}

func (p *Provider) ProvideSubscriptions(database *sql_db.SqlDatabase) *[]bus.EventSubscription {
	return &[]bus.EventSubscription{}
}

func (p *Provider) ProvideKafkaConsumer(eventBus *bus.EventBus) (*kafka.KafkaConsumer, error) {
	brokers := p.kafkaBrokers()

	return kafka.NewKafkaConsumer(brokers, eventBus)
}

func (p *Provider) kafkaBrokers() []string {
	if p.env == "development" {
		return []string{
			"localhost:9093",
		}
	} else {
		return []string{
			"172.31.36.175:9092",
			"172.31.45.255:9092",
		}
	}
}
