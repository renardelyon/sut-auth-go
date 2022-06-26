package main

import (
	"fmt"
	"log"
	"net"
	"sut-auth-go/config"
	service "sut-auth-go/domain/auth"
	db "sut-auth-go/lib/pkg"
	"sut-auth-go/lib/utils"
	pb "sut-auth-go/pb/auth"

	"google.golang.org/grpc"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed at config: ", err.Error())
	}

	h := db.Init(c.DBUrl)

	jwt := utils.JwtWrapper{
		SecretKey:       c.JWTKey,
		Issuer:          "sut-auth-go",
		ExpirationHours: 1,
	}

	lis, err := net.Listen("tcp", c.Port)
	if err != nil {
		log.Fatalf("Cannot listen PORT: %s", c.Port)
	}

	fmt.Println("Auth service on Port: ", c.Port)

	s := service.Service{
		H:   h,
		Jwt: jwt,
		C:   c,
	}

	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}

}
