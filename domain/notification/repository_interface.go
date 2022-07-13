package notification

import notifpb "sut-auth-go/pb/notification"

type NotificationRepoInterface interface {
	SubscribeNotificationByUserId(userId string) (*notifpb.SubscribeNotificationResponse, error)
}
