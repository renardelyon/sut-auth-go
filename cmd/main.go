package main

import (
	"log"
	"sut-auth-go/application"
	"sut-auth-go/config"
	"sut-auth-go/domain/auth/service"
	pb "sut-auth-go/pb/auth"
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

	s := service.NewService(app.DbClients, app.JwtWrapper, c)

	pb.RegisterAuthServiceServer(app.GrpcServer, s)

	err = app.Run(&c)
	if err != nil {
		log.Fatalln(err.Error())
	}

}
