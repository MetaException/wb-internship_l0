package natsbroker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func NewConsumer(ctx context.Context, nc *nats.Conn) (jetstream.Consumer, error) {

	js, err := jetstream.New(nc)
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	c, err := js.CreateOrUpdateConsumer(ctx, "L0_STREAM", jetstream.ConsumerConfig{
		Durable:       "CONS",
		AckPolicy:     jetstream.AckExplicitPolicy,
		FilterSubject: "l0.*",
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create consumer for stream L0_STREAM: %w", err)
	}

	return c, nil
}

func (ns NatsBroker) ConsumeOrders() func(msg jetstream.Msg) {
	return func(msg jetstream.Msg) {
		msg.Ack()

		msg_data := msg.Data()

		fmt.Printf("Received a JetStream message via callback: %s\n", string(msg_data))

		var orderinfo map[string]interface{}
		if err := json.Unmarshal(msg_data, &orderinfo); err != nil {
			fmt.Printf("Error unmarshalling JSON for UUID: %v\n", err)
		}

		ns.Cache.Set(orderinfo["order_uid"].(string), msg_data)
		ns.DB.AddToDb(orderinfo["order_uid"].(string), msg_data)

		fmt.Println(orderinfo["order_uid"].(string)) //TODO упроситт
	}
}
