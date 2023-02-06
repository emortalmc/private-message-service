package notifier

import (
	"context"
	"fmt"
	"github.com/emortalmc/proto-specs/gen/go/message/privatemessage"
	pb "github.com/emortalmc/proto-specs/gen/go/model/privatemessage"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
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

func (n *rabbitMqNotifier) MessageSent(ctx context.Context, pm *pb.PrivateMessage) error {
	message := &privatemessage.PrivateMessageReceivedMessage{PrivateMessage: pm}

	bytes, err := proto.Marshal(message)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 5)
	defer cancel()

	return n.channel.PublishWithContext(ctx, "mc:proxy:all", "", false, false, amqp.Publishing{
		ContentType: "application/x-protobuf",
		Type:        string(message.ProtoReflect().Descriptor().FullName()),
		Body:        bytes,
	})
}
