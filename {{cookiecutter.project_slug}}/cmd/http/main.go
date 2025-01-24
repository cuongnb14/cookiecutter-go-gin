package main

import (
	"context"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"{{ cookiecutter.project_slug }}/configs"
	"{{ cookiecutter.project_slug }}/internal/controllers"
	"{{ cookiecutter.project_slug }}/internal/core/repositories"
	"{{ cookiecutter.project_slug }}/internal/core/services"
	"{{ cookiecutter.project_slug }}/internal/server"
)

func main() {
	configs.PreServerStart()
	defer configs.PreServerShutdown()

	defer sentry.Flush(2 * time.Second)

	configs.InitAsynqClient()
	defer configs.AsynqClient.Close()

	Run()

}

func Run() {
	app := fx.New(
		fx.Provide(
			configs.GetDB,
			repositories.NewUserRepository,
			services.NewUserService,
			controllers.NewUserController,
			server.NewRouter,
			server.NewServer,
		),
		fx.Invoke(startServer),
	)

	app.Run()
}

func startServer(lc fx.Lifecycle, s *server.Server) {
	logger := configs.GetLogger()

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				gin.SetMode(configs.Env.GinMode)
				s.Run()
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("wait for server stop")
			if err := s.HttpServer.Shutdown(ctx); err != nil {
				logger.Error("server forced to shutdown:", "err", err)
			}
			return nil
		},
	})
}
