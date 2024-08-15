package controllers

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"{{ cookiecutter.project_slug }}/configs"
	"{{ cookiecutter.project_slug }}/internal/core/vo"

	"github.com/gin-gonic/gin"
)

func UserLoginHandler(ctx *gin.Context) {
	data := vo.UserVO{}
	responses.Ok(ctx, data)
}
