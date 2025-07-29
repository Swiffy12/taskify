package kafkaclient

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	log  *slog.Logger
	conn *kafka.Conn
}

func New(
	ctx context.Context,
	log *slog.Logger,
	addr string,
	topic string,
	deadline time.Duration,
	partition int,
) *Producer {
	conn, err := kafka.DialLeader(ctx, "tcp", addr, topic, partition)
	if err != nil {
		fmt.Println(err)
		panic("failed to connect to kafka")
	}

	return &Producer{conn: conn, log: log}
}

func (p *Producer) Send(msg string) error {
	const op = "producer.Send"

	_, err := p.conn.WriteMessages(kafka.Message{Value: []byte(msg)})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
