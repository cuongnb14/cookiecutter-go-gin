package server

import (
	"net/http"
	"time"

	"{{ cookiecutter.project_slug }}/configs"
	"{{ cookiecutter.project_slug }}/internal/controllers"
	"{{ cookiecutter.project_slug }}/internal/core/repositories"
	"{{ cookiecutter.project_slug }}/internal/core/services"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Router     *gin.Engine
	HttpServer *http.Server
}

func (server *Server) Initialize() {
	gin.SetMode(configs.Env.GinMode)

	userRepo := repositories.NewUserRepository(configs.GetDB())
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	server.Router = NewRouter(userController)
	// server.Router.Use(gin.Recovery(), middlewares.Logger())

}

func (server *Server) Run() {
	logger := configs.GetLogger()
	logger.Info("Listening to address: 0.0.0.0:" + configs.Env.Port)
	server.HttpServer = &http.Server{
		Addr:         ":" + configs.Env.Port,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		// MaxHeaderBytes: 1 << 20,
		Handler: server.Router,
	}
	server.HttpServer.ListenAndServe()
	// server.Router.Run()
}
