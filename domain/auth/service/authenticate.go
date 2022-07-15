package service

import (
	"context"
	"net/http"
	"sut-auth-go/domain/auth/model"
	pb "sut-auth-go/pb/auth"
)

func (s *Service) Authenticate(ctx context.Context, reqAuth *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	claims, err := s.Jwt.AuthenticateToken(reqAuth.Token)

	if err != nil {
		return &pb.AuthenticateResponse{
			Status: http.StatusUnauthorized,
			Error:  err.Error(),
		}, nil
	}

	var user model.User
	var admin model.Admin

	if result := s.H.DB.Where(&model.User{Id: claims.Id}).First(&user); result.Error != nil {
		if result := s.H.DB.Where(&model.Admin{Id: claims.Id}).First(&admin); result.Error != nil {
			return &pb.AuthenticateResponse{
				Status: http.StatusNotFound,
				Error:  "User Not Found",
			}, nil
		}
	}

	var auth *pb.UserInfo

	if user.Id != "" {
		auth = &pb.UserInfo{
			Id:       user.Id,
			Username: user.Username,
			Name:     user.Name,
			Role:     pb.Role_USER,
			AdminId:  user.AdminId,
		}
	} else {
		auth = &pb.UserInfo{
			Id:       admin.Id,
			Username: admin.Username,
			Name:     admin.Name,
			Role:     pb.Role_ADMIN,
		}
	}

	return &pb.AuthenticateResponse{
		Status:   http.StatusOK,
		Userinfo: auth,
	}, nil

}
