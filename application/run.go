package application

import "sut-auth-go/config"

func (app *Application) Run(cfg *config.Config) error {
	return grpcRun(cfg)(app)
}
