package server

import (
	"net/http"
	"time"

	"{{ cookiecutter.project_slug }}/configs"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Router     *gin.Engine
	HttpServer *http.Server
}

func NewServer(router *gin.Engine) *Server {
	return &Server{
		Router: router,
	}
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
