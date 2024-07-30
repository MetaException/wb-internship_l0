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
}

func New(cacheStorage *cachestorage.CacheStorage, pg *postgresql.Postgres, logger *logrus.Logger) *NatsBroker {

	nc, err := nats.Connect("nats://localhost:4223")

	if err != nil {
		logrus.WithError(err).Fatal("failed to connect to NATS")
	}

	return &NatsBroker{
		Logger:       logger,
		CacheStorage: cacheStorage,
		DB:           pg,
		conn:         nc,
	}
}

func (ns *NatsBroker) Close() {
	ns.conn.Close()
}

func (ns *NatsBroker) Listen() error {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	c, err := NewConsumer(ctx, ns.conn)

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
