package notifier

import (
	"context"
	"fmt"
	pb "github.com/emortalmc/proto-specs/gen/go/model/privatemessage"
	"github.com/golang/protobuf/proto"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"private-message-service/internal/config"
)

const rabbitMqUriFormat = "amqp://%s:%s@%s:5672/"

type rabbitMqNotifier struct {
	Notifier
	channel *amqp.Channel
}

func NewRabbitMqNotifier(cfg config.RabbitMQConfig) (Notifier, error) {
	conn, err := amqp.Dial(fmt.Sprintf(rabbitMqUriFormat, cfg.Username, cfg.Password, cfg.Host))
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &rabbitMqNotifier{
		channel: channel,
	}, nil
}

func (n *rabbitMqNotifier) MessageSent(ctx context.Context, msg *pb.PrivateMessage) error {
	bytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 5)
	defer cancel()

	log.Print(msg.ProtoReflect().Descriptor().FullName())
	return n.channel.PublishWithContext(ctx, "mc:proxy:all", "", false, false, amqp.Publishing{
		ContentType: "application/x-protobuf",
		Type:        string(msg.ProtoReflect().Descriptor().FullName()),
		Body:        bytes,
	})
}
