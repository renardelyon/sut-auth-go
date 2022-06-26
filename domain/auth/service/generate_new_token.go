package service

import (
	"context"
	"net/http"
	"sut-auth-go/domain/auth/model"
	"sut-auth-go/lib/utils"
	pb "sut-auth-go/pb/auth"
)

func (s *Service) GenerateNewToken(ctx context.Context, reqToken *pb.GenerateNewTokenRequest) (*pb.GenerateNewTokenResponse, error) {
	refreshToken := reqToken.RefreshToken
	var tokens model.Token

	if refreshToken == "" {
		return &pb.GenerateNewTokenResponse{
			Status: http.StatusBadRequest,
			Error:  "refresh token request body is null",
		}, nil
	}

	if result := s.H.DB.Where(&model.Token{Token: refreshToken}).First(&tokens); result.Error != nil {
		return &pb.GenerateNewTokenResponse{
			Status: http.StatusNotFound,
			Error:  "refresh token not found in the database",
		}, nil
	}

	newJwt := utils.JwtWrapper{
		SecretKey:       s.Jwt.SecretKey,
		Issuer:          s.Jwt.Issuer,
		ExpirationHours: 24 * 7,
	}

	claims, err := s.Jwt.AuthenticateToken(refreshToken)

	if err != nil {
		return &pb.GenerateNewTokenResponse{
			Status: http.StatusUnauthorized,
			Error:  "token can't be authenticated",
		}, nil
	}

	newToken, _ := newJwt.GenerateAccessToken(utils.UserAdmin{
		Id:       claims.Id,
		Username: claims.Username,
		Name:     claims.Name,
		Role:     claims.Role,
	})

	return &pb.GenerateNewTokenResponse{
		Status:   http.StatusOK,
		NewToken: newToken,
	}, nil
}
