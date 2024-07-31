package natsbroker

import (
	"context"
	"time"

	"github.com/MetaException/wb_l0/internal/cachestorage"
	"github.com/MetaException/wb_l0/internal/postgresql"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type NatsBroker struct {
	Logger       *logrus.Logger
	CacheStorage *cachestorage.CacheStorage
	DB           *postgresql.Postgres
	conn         *nats.Conn
	config       *Config
}

func New(cacheStorage *cachestorage.CacheStorage, pg *postgresql.Postgres, logger *logrus.Logger, config *Config) *NatsBroker {

	nc, err := nats.Connect(config.URL)

	if err != nil {
		logger.WithError(err).Fatal("failed to connect to NATS")
	}

	logger.Info("successfully connected to NATS")

	return &NatsBroker{
		Logger:       logger,
		CacheStorage: cacheStorage,
		DB:           pg,
		conn:         nc,
		config:       config,
	}
}

func (ns *NatsBroker) Close() {
	ns.conn.Close()
}

func (ns *NatsBroker) Listen() error {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	c, err := ns.NewConsumer(ctx)

	if err != nil {
		ns.Logger.WithError(err).Error("failed to create consumer")
		return errors.WithStack(err)
	}

	_, err = c.Consume(ns.ConsumeOrders())

	if err != nil {
		ns.Logger.WithError(err).Error("NATS consume error")
		return errors.WithStack(err)
	}

	return nil
}
