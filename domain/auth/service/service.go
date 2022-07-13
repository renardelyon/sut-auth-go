package service

import (
	"sut-auth-go/config"
	"sut-auth-go/domain/notification"
	"sut-auth-go/lib/pkg/db"
	"sut-auth-go/lib/utils"
)

type Service struct {
	H              db.Handler
	Jwt            utils.JwtWrapper
	C              config.Config
	NotifInterface notification.NotificationRepoInterface
}

func NewService(h db.Handler, jwt utils.JwtWrapper, c config.Config, notifInterface notification.NotificationRepoInterface) *Service {
	return &Service{
		H:              h,
		Jwt:            jwt,
		C:              c,
		NotifInterface: notifInterface,
	}
}
