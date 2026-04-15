package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Handler func(ctx context.Context, event OrderEvent) error

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

func New(url, queue string) (*Consumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("channel: %w", err)
	}

	_, err = ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("declare queue: %w", err)
	}

	// не брать больше одного сообщения пока не обработано текущее
	ch.Qos(1, 0, false)

	return &Consumer{conn: conn, channel: ch, queue: queue}, nil
}

func (c *Consumer) Run(ctx context.Context, handler Handler) error {
	msgs, err := c.channel.Consume(c.queue, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("consume: %w", err)
	}

	log.Printf("waiting for messages on queue %q", c.queue)

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-msgs:
			if !ok {
				return fmt.Errorf("channel closed")
			}

			var event OrderEvent
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Printf("bad message, skipping: %v", err)
				msg.Nack(false, false)
				continue
			}

			if err := handler(ctx, event); err != nil {
				log.Printf("handler error order_id=%d: %v", event.OrderID, err)
				msg.Nack(false, true) // вернуть в очередь
				continue
			}

			msg.Ack(false)
		}
	}
}

func (c *Consumer) Close() {
	c.channel.Close()
	c.conn.Close()
}
