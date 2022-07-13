package application

import (
	"context"
	"log"
	"net"
	"sut-auth-go/config"
	"sut-auth-go/lib/pkg/db"
	"sut-auth-go/lib/utils"

	"google.golang.org/grpc"
)

func initGrpcServer(cfg *config.Config) func(*Application) error {
	return func(app *Application) error {
		g := grpc.NewServer()
		app.GrpcServer = g
		return nil
	}
}

func grpcRun(cfg *config.Config) func(*Application) error {
	return func(app *Application) error {
		lis, err := net.Listen("tcp", cfg.Port)
		if err != nil {
			return err
		}
		log.Println("Auth service on Port: ", cfg.Port)
		if err := app.GrpcServer.Serve(lis); err != nil {
			return err
		}
		app.GrpcServer.GracefulStop()
		return nil
	}
}

func Setup(cfg *config.Config) (*Application, error) {
	app := new(Application)
	err := runInit(
		initDatabase(cfg),
		initGrpcClient(cfg),
		initApp(cfg),
		initJwtWrapper(cfg))(app)

	if err != nil {
		return app, err
	}
	return app, nil
}

func runInit(appFuncs ...func(*Application) error) func(*Application) error {
	return func(app *Application) error {
		app.Context = context.Background()
		for _, fn := range appFuncs {
			if e := fn(app); e != nil {
				return e
			}
		}
		return nil
	}
}

func initGrpcClient(cfg *config.Config) func(*Application) error {
	return func(app *Application) error {
		app.GrpcClients = make(map[string]*grpc.ClientConn)

		NotifServiceCfg := cfg.NotifHost
		conn, err := setupGrpcConnection(NotifServiceCfg)
		if err != nil {
			return err
		}

		app.GrpcClients["notification-service"] = conn
		log.Println("notification-service" + " connected on " + NotifServiceCfg)

		log.Println("init Grpc Client done")
		return nil
	}
}

func setupGrpcConnection(cfg string) (*grpc.ClientConn, error) {
	return grpc.Dial(cfg, grpc.WithInsecure())
}

func initJwtWrapper(cfg *config.Config) func(*Application) error {
	return func(app *Application) error {
		app.JwtWrapper = utils.JwtWrapper{
			SecretKey:       cfg.JWTKey,
			Issuer:          "sut-auth-go",
			ExpirationHours: 1,
		}

		log.Println("init jwt wrapper done")
		return nil
	}
}

func initDatabase(cfg *config.Config) func(*Application) error {
	return func(app *Application) error {
		app.DbClients = db.Init(cfg.DBUrl)

		log.Println("init postgre database done")

		return nil
	}
}

func initApp(cfg *config.Config) func(*Application) error {
	return func(app *Application) error {
		return initGrpcServer(cfg)(app)
	}
}
