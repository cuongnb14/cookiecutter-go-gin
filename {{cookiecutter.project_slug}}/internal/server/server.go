package server

import (
	"net/http"
	"time"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"{{ cookiecutter.project_slug }}/configs"
	"{{ cookiecutter.project_slug }}/internal/middlewares"
)

type Server struct {
	Router     *gin.Engine
	HttpServer *http.Server
}

func (server *Server) Initialize() {
	gin.SetMode(configs.Env.GinMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}))
	router.Use(middlewares.ErrorHandler())
	router.Use(middlewares.RequestLogger())

	server.Router = router
	InitRouteV1(server.Router)
	// server.Router.Use(gin.Recovery(), middlewares.Logger())

}

func (server *Server) Run() {
	logger := configs.GetLogger()
	logger.Info("Listening to port" + configs.Env.Port.String())
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
