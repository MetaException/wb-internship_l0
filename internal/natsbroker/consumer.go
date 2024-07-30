package natsbroker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func NewConsumer(ctx context.Context, nc *nats.Conn) (jetstream.Consumer, error) {

	js, err := jetstream.New(nc)
	if err != nil {
		logrus.WithError(err).Error("failed to create JetStream context")
		return nil, errors.WithStack(err)
	}

	c, err := js.CreateOrUpdateConsumer(ctx, "L0_STREAM", jetstream.ConsumerConfig{
		Durable:       "CONS",
		AckPolicy:     jetstream.AckExplicitPolicy,
		FilterSubject: "l0.*",
	})

	if err != nil {
		logrus.WithError(err).Error("failed to create consumer for stream L0_STREAM")
		return nil, errors.WithStack(err)
	}

	return c, nil
}

func (ns *NatsBroker) ConsumeOrders() func(msg jetstream.Msg) {
	return func(msg jetstream.Msg) {
		msg.Ack()

		msg_data := msg.Data()

		ns.Logger.Infof("received a JetStream message via callback: %s", string(msg_data))

		var orderinfo map[string]interface{}
		if err := json.Unmarshal(msg_data, &orderinfo); err != nil {
			ns.Logger.WithError(err).Error("error unmarshalling JSON for UUID")
		}

		ns.CacheStorage.Cache.Set(orderinfo["order_uid"].(string), msg_data)
		ns.DB.AddToDb(orderinfo["order_uid"].(string), msg_data)

		fmt.Println(orderinfo["order_uid"].(string)) //TODO упроситт
	}
}
