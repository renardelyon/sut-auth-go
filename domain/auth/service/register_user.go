package service

import (
	"context"
	"log"
	"net/http"
	"sut-auth-go/domain/auth/model"
	"sut-auth-go/lib/utils"
	pb "sut-auth-go/pb/auth"

	"github.com/google/uuid"
)

func (s *Service) RegisterUser(ctx context.Context, reqUser *pb.UserRegisterRequest) (*pb.UserRegisterResponse, error) {
	var user model.User
	if result := s.H.DB.Where(&model.User{Username: reqUser.Username}).First(&user); result.Error == nil {
		return &pb.UserRegisterResponse{
			Status: http.StatusConflict,
			Error:  "Username Already Exists",
		}, nil
	}

	uuid, _ := uuid.NewUUID()

	newUser := &model.User{
		Id:       uuid.String(),
		Username: reqUser.Username,
		AdminId:  reqUser.AdminId,
		Name:     reqUser.Name,
		Role:     "user",
		Password: utils.HashPassword(reqUser.Password),
	}

	if res := s.H.DB.Create(newUser); res.Error != nil {
		log.Println(res.Error)
		return &pb.UserRegisterResponse{
			Status: http.StatusInternalServerError,
			Error:  res.Error.Error(),
		}, nil
	}

	resp, _ := s.NotifInterface.SubscribeNotificationByUserId(uuid.String())
	if resp.Error != "" {
		log.Println(resp.Error)
		return &pb.UserRegisterResponse{
			Status: http.StatusInternalServerError,
			Error:  resp.Error,
		}, nil
	}

	return &pb.UserRegisterResponse{
		Status: http.StatusCreated,
	}, nil
}
