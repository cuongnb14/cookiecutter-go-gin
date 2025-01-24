package server

import (
	"time"

	"{{ cookiecutter.project_slug }}/internal/controllers"
	"{{ cookiecutter.project_slug }}/internal/middlewares"

	"net/http"

	// _ "{{ cookiecutter.project_slug }}/docs"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	gin "github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(userController *controllers.UserController) *gin.Engine {
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

	router.GET("/login", userController.SignUp)

	// api := r.Group("/api/wallet/v1")
	// api.POST("/login", controllers.PlayerLoginHandler)

	// 	api.POST("/users/:id", middlewares.AuthMiddleware(), s.DeleteUser)

	// Docs
	router.GET("docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Not Found Route
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "ErrNotFound", "error": "ErrNotFound"})
	})

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	return router
}
