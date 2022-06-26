package service

import (
	"context"
	"net/http"
	"sut-auth-go/config"
	"sut-auth-go/domain/auth/model"
	db "sut-auth-go/lib/pkg"
	"sut-auth-go/lib/utils"
	pb "sut-auth-go/pb/auth"

	"github.com/google/uuid"
)

type Service struct {
	H   db.Handler
	Jwt utils.JwtWrapper
	C   config.Config
}

func (s *Service) RegisterUser(ctx context.Context, reqUser *pb.UserRegisterRequest) (*pb.UserRegisterResponse, error) {
	var user model.User
	if result := s.H.DB.Where(&model.User{Username: reqUser.Username}).First(&user); result.Error == nil {
		return &pb.UserRegisterResponse{
			Status: http.StatusConflict,
			Error:  "Email Already Exists",
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

	s.H.DB.Create(newUser)

	// TODO: connect to notification service

	return &pb.UserRegisterResponse{
		Status: http.StatusCreated,
	}, nil
}

func (s *Service) RegisterAdmin(ctx context.Context, reqAdmin *pb.AdminRegisterRequest) (*pb.AdminRegisterResponse, error) {
	var admin model.Admin
	if result := s.H.DB.Where(&model.Admin{Username: reqAdmin.Username}).First(&admin); result.Error == nil {
		return &pb.AdminRegisterResponse{
			Status: http.StatusConflict,
			Error:  "Email Already Exists",
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

	return &pb.LoginResponse{
		Status: http.StatusOK,
		Token:  token,
	}, nil
}

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
