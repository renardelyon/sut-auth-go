package main

import (
	"log"
	"sut-auth-go/application"
	"sut-auth-go/config"
	"sut-auth-go/domain/auth/service"
	notifGrpc "sut-auth-go/domain/notification/repo/grpc"
	pb "sut-auth-go/pb/auth"
	notifpb "sut-auth-go/pb/notification"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed at config: ", err.Error())
	}

	app, err := application.Setup(&c)
	if err != nil {
		log.Fatalln("Failed at application setup: ", err.Error())
	}

	notifClient := notifpb.NewNotificationServiceClient(app.GrpcClients["notification-service"])
	notifRepo := notifGrpc.NewGrpcRepo(notifClient)

	s := service.NewService(app.DbClients, app.JwtWrapper, c, notifRepo)

	pb.RegisterAuthServiceServer(app.GrpcServer, s)

	err = app.Run(&c)
	if err != nil {
		log.Fatalln(err.Error())
	}

}
