package service

import (
	"context"
	"net/http"
	"sut-auth-go/domain/auth/model"
	pb "sut-auth-go/pb/auth"
)

func (s *Service) DeleteToken(ctx context.Context, reqDel *pb.DeleteTokenRequest) (*pb.DeleteTokenResponse, error) {
	if result := s.H.DB.Where("token = ?", reqDel.Token).Delete(&model.Token{}); result.Error != nil {
		return &pb.DeleteTokenResponse{
			Status: http.StatusNotFound,
			Error:  result.Error.Error(),
		}, nil
	}

	return &pb.DeleteTokenResponse{
		Status: http.StatusOK,
	}, nil
}
