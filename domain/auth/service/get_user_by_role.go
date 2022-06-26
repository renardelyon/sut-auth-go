package service

import (
	"context"
	"net/http"
	"sut-auth-go/domain/auth/model"
	pb "sut-auth-go/pb/auth"
)

func (s *Service) GetUserByRole(ctx context.Context, reqUa *pb.GetUserByRoleRequest) (*pb.GetUserByRoleResponse, error) {
	role := reqUa.Role
	var user []model.User
	var admin []model.Admin

	var infoArray []*pb.UserAdminInfo

	s.H.DB.Where(&model.User{Role: role}).Find(&user)
	s.H.DB.Where(&model.Admin{Role: role}).Find(&admin)
	if len(user) <= 0 {
		if len(admin) <= 0 {
			return &pb.GetUserByRoleResponse{
				Status: http.StatusNotFound,
				Error:  "role doesn't exist",
			}, nil
		}
	}

	if len(user) > 0 {
		for _, each := range user {
			infoArray = append(infoArray, &pb.UserAdminInfo{
				Id:       each.Id,
				Username: each.Username,
				Name:     each.Name,
				Role:     pb.Role_USER.String(),
			})
		}
	}

	if len(admin) > 0 {
		for _, each := range admin {
			infoArray = append(infoArray, &pb.UserAdminInfo{
				Id:       each.Id,
				Username: each.Username,
				Name:     each.Name,
				Role:     pb.Role_ADMIN.String(),
			})
		}
	}

	return &pb.GetUserByRoleResponse{
		Status:        http.StatusOK,
		UserAdminInfo: infoArray,
	}, nil
}
