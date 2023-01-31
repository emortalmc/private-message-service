package app

import (
	"context"
	"fmt"
	"github.com/emortalmc/proto-specs/gen/go/grpc/privatemessage"
	"github.com/emortalmc/proto-specs/gen/go/grpc/relationship"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"private-message-service/internal/config"
	"private-message-service/internal/notifier"
	"private-message-service/internal/service"
)

func Run(ctx context.Context, cfg *config.Config, logger *zap.SugaredLogger) {
	notif, err := notifier.NewRabbitMqNotifier(cfg.RabbitMQ)
	if err != nil {
		logger.Fatalw("failed to create notifier", "error", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		logger.Fatalw("failed to listen", "error", err)
	}

	conn, err := grpc.Dial(cfg.RelationshipService.Host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalw("failed to connect to relationship service", "error", err)
	}
	rc := relationship.NewRelationshipClient(conn)

	s := grpc.NewServer()
	privatemessage.RegisterPrivateMessageServer(s, service.NewPrivateMessageService(notif, rc))
	logger.Infow("listening on port", "port", cfg.Port)

	err = s.Serve(lis)
	if err != nil {
		logger.Fatalw("failed to serve", "error", err)
		return
	}
}