package natsbroker

import (
	"context"
	"encoding/json"

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

		if !validate(orderinfo) {
			ns.Logger.Error("invalid data received")
			return
		}

		orderUID, ok := orderinfo["order_uid"].(string)
		if !ok {
			ns.Logger.Error("Invalid or missing 'order_uid'")
			return
		}

		ns.CacheStorage.Cache.Set(orderUID, msg_data)
		if err := ns.DB.AddToDb(orderUID, msg_data); err != nil {
			ns.Logger.WithError(err).Error("Error saving to database")
			return
		}
	}
}

func validate(data map[string]interface{}) bool {
	return data["order_uid"] != nil
}
