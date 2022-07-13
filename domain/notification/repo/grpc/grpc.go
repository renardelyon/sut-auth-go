package grpc

import (
	"context"
	notifpb "sut-auth-go/pb/notification"
)

type repo struct {
	notifGrpc notifpb.NotificationServiceClient
}

func NewGrpcRepo(notifGrpc notifpb.NotificationServiceClient) *repo {
	return &repo{
		notifGrpc: notifGrpc,
	}
}

func (r *repo) SubscribeNotificationByUserId(userId string) (*notifpb.SubscribeNotificationResponse, error) {
	req := &notifpb.SubscribeNotificationRequest{
		UserId: userId,
	}
	return r.notifGrpc.SubscribeNotificationByUserId(context.Background(), req)
}
