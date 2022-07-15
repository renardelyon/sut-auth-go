package service

import (
	"context"
	"net/http"
	"strings"
	"sut-auth-go/domain/auth/model"
	"sut-auth-go/lib/utils"
	pb "sut-auth-go/pb/auth"
)

func (s *Service) Login(ctx context.Context, reqLogin *pb.LoginRequest) (*pb.LoginResponse, error) {
	var user model.User
	var admin model.Admin
	var ua utils.UserAdmin

	if result := s.H.DB.Where(&model.User{Username: reqLogin.Username}).First(&user); result.Error != nil {
		if result := s.H.DB.Where(&model.Admin{Username: reqLogin.Username}).First(&admin); result.Error != nil {
			return &pb.LoginResponse{
				Status: http.StatusNotFound,
				Error:  "User Not Found",
				Token:  "",
			}, nil
		}
	}

	var match bool

	if user.Id != "" {
		match = utils.CompareHashPassword(reqLogin.Password, user.Password)
		ua = utils.UserAdmin{
			Id:       user.Id,
			Username: user.Username,
			Name:     user.Name,
			Role:     user.Role,
		}
	} else {
		match = utils.CompareHashPassword(reqLogin.Password, admin.Password)
		ua = utils.UserAdmin{
			Id:       admin.Id,
			Username: admin.Username,
			Name:     admin.Name,
			Role:     admin.Role,
		}
	}

	if !match {
		return &pb.LoginResponse{
			Status: http.StatusUnauthorized,
			Error:  "Password incorrect",
			Token:  "",
		}, nil
	}

	token, err := s.Jwt.GenerateAccessToken(ua)
	if err != nil {
		return &pb.LoginResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}, nil
	}

	newJwt := utils.JwtWrapper{
		SecretKey:       s.Jwt.SecretKey,
		Issuer:          s.Jwt.Issuer,
		ExpirationHours: 24 * 365,
	}

	refreshToken, err := newJwt.GenerateAccessToken(ua)
	if err != nil {
		return &pb.LoginResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}, nil
	}

	s.H.DB.Create(&model.Token{
		Token: refreshToken,
	})

	return &pb.LoginResponse{
		Status:       http.StatusOK,
		Token:        token,
		Refreshtoken: refreshToken,
		UserInfo: &pb.UserInfo{
			Id:       ua.Id,
			Username: ua.Name,
			Name:     ua.Name,
			Role:     pb.Role(pb.Role_value[strings.ToUpper(ua.Role)]),
		},
	}, nil
}
