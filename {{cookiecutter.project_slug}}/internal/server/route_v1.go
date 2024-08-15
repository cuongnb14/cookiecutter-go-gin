package server

import (
	"net/http"

	_ "{{ cookiecutter.project_slug }}/docs"

	gin "github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouteV1(r *gin.Engine) {

	// api := r.Group("/api/wallet/v1")
	// api.POST("/login", controllers.PlayerLoginHandler)

	// 	api.POST("/users/:id", middlewares.AuthMiddleware(), s.DeleteUser)

	// Docs
	r.GET("docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Not Found Route
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "ErrNotFound", "error": "ErrNotFound"})
	})

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

}
