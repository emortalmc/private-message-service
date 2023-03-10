package service

import (
	"context"
	pb "github.com/emortalmc/proto-specs/gen/go/grpc/privatemessage"
	"github.com/emortalmc/proto-specs/gen/go/grpc/relationship"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"private-message-service/internal/notifier"
)

type privateMessageService struct {
	pb.PrivateMessageServer
	notif notifier.Notifier
	rs    relationship.RelationshipClient
}

func NewPrivateMessageService(notif notifier.Notifier, rs relationship.RelationshipClient) pb.PrivateMessageServer {
	return &privateMessageService{
		notif: notif,
		rs:    rs,
	}
}

func (s *privateMessageService) SendPrivateMessage(ctx context.Context, req *pb.PrivateMessageRequest) (*pb.PrivateMessageResponse, error) {
	// TODO: Check if the player is online.

	resp, err := s.rs.IsBlocked(ctx, &relationship.IsBlockedRequest{
		IssuerId: req.Message.SenderId,
		TargetId: req.Message.RecipientId,
	})
	if err != nil {
		return nil, err
	}

	block := resp.GetBlock()

	if block != nil {
		if block.BlockerId == req.Message.SenderId {
			st := status.New(codes.PermissionDenied, "you have blocked this player")
			st, _ = st.WithDetails(&pb.PrivateMessageErrorResponse{Reason: pb.PrivateMessageErrorResponse_YOU_BLOCKED})
			return nil, st.Err()
		} else {
			st := status.New(codes.PermissionDenied, "you are blocked by this player")
			st, _ = st.WithDetails(&pb.PrivateMessageErrorResponse{Reason: pb.PrivateMessageErrorResponse_PRIVACY_BLOCKED})
			return nil, st.Err()
		}
	}

	err = s.notif.MessageSent(ctx, req.Message)
	if err != nil {
		return nil, err
	}

	return &pb.PrivateMessageResponse{
		Message: req.Message,
	}, nil
}
