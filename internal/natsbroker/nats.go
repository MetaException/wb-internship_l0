package natsbroker

import (
	"context"
	"fmt"
	"time"

	"github.com/MetaException/wb_l0/internal/cache"
	"github.com/MetaException/wb_l0/internal/postgresql"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

type NatsBroker struct {
	Logger *logrus.Logger
	Cache  *cache.CacheStorage
	DB     *postgresql.Postgres
	conn   *nats.Conn
}

func New(cacheStorage *cache.CacheStorage, pg *postgresql.Postgres) *NatsBroker {

	nc, err := nats.Connect("nats://localhost:4223")

	if err != nil {
		logrus.Fatal("failed to connect to NATS: %w", err)
	}

	return &NatsBroker{
		Logger: logrus.New(),
		Cache:  cacheStorage,
		DB:     pg,
		conn:   nc,
	}
}

func (ns NatsBroker) Listen() error {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	c, err := NewConsumer(ctx, ns.conn)

	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}

	_, err = c.Consume(ns.ConsumeOrders())

	if err != nil {
		return fmt.Errorf("consume error: %w", err)
	}

	return nil
}
