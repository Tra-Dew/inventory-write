package cmd

import (
	"context"
	"fmt"

	"github.com/Tra-Dew/inventory-write/pkg/core"
	"github.com/Tra-Dew/inventory-write/pkg/inventory"
	"github.com/Tra-Dew/inventory-write/pkg/inventory/postgres"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

// Container contains all depencies from our api
type Container struct {
	Settings *core.Settings

	DBConnPool *pgxpool.Pool

	Authenticate *core.Authenticate

	Producer *core.MessageBrokerProducer
	SNS      *session.Session
	SQS      *session.Session

	InventoryRepository inventory.Repository
	InventoryService    inventory.Service
	InventoryController inventory.Controller
}

// NewContainer creates new instace of Container
func NewContainer(settings *core.Settings) *Container {

	container := new(Container)

	container.Settings = settings

	container.DBConnPool = connectPostgres(settings.Postgres)

	container.SQS = session.Must(session.NewSession(&aws.Config{
		Region:   aws.String(settings.SQS.Region),
		Endpoint: aws.String(settings.SQS.Endpoint),
	}))

	container.SNS = session.Must(session.NewSession(&aws.Config{
		Region:   aws.String(settings.SNS.Region),
		Endpoint: aws.String(settings.SNS.Endpoint),
	}))

	container.Producer = core.NewMessageBrokerProducer(container.SNS)

	container.Authenticate = core.NewAuthenticate(settings.JWT.Secret)

	container.InventoryRepository = postgres.NewRepository(container.DBConnPool)
	container.InventoryService = inventory.NewService(container.InventoryRepository)
	container.InventoryController = inventory.NewController(settings, container.Authenticate, container.InventoryService)

	return container
}

// Controllers maps all routes and exposes them
func (c *Container) Controllers() []core.Controller {
	return []core.Controller{
		&c.InventoryController,
	}
}

// Close terminates every opened resource
func (c *Container) Close() {}

func connectPostgres(conf *core.PostgresConfig) *pgxpool.Pool {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", conf.User, conf.Password, conf.Host, conf.Port, conf.Database)

	pool, err := pgxpool.Connect(context.Background(), connString)

	if err != nil {
		logrus.
			WithError(err).
			Fatalf("unable to connect to database")
	}

	if err = pool.Ping(context.Background()); err != nil {
		logrus.
			WithError(err).
			Fatalf("unable to ping database")
	}

	logrus.Info("connected to postgres")

	return pool
}
