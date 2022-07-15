package service

import (
	"context"
	"net/http"
	"sut-auth-go/domain/auth/model"
	"sut-auth-go/lib/utils"
	pb "sut-auth-go/pb/auth"

	"github.com/google/uuid"
)

func (s *Service) RegisterAdmin(ctx context.Context, reqAdmin *pb.AdminRegisterRequest) (*pb.AdminRegisterResponse, error) {
	var admin model.Admin
	if result := s.H.DB.Where(&model.Admin{Username: reqAdmin.Username}).First(&admin); result.Error == nil {
		return &pb.AdminRegisterResponse{
			Status: http.StatusConflict,
			Error:  "Username Already Exists",
		}, nil
	}

	uuid, _ := uuid.NewUUID()

	newAdmin := &model.Admin{
		Id:       uuid.String(),
		Username: reqAdmin.Username,
		Name:     reqAdmin.Name,
		Role:     "admin",
		Password: utils.HashPassword(reqAdmin.Password),
	}

	if s.C.AdminKey == reqAdmin.AdminKey {
		s.H.DB.Create(newAdmin)
	} else {
		return &pb.AdminRegisterResponse{
			Status: http.StatusBadRequest,
			Error:  "Admin key isn't valid",
		}, nil
	}

	return &pb.AdminRegisterResponse{
		Status: http.StatusCreated,
	}, nil
}
