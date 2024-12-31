package controllers

import (
	"{{ cookiecutter.project_slug }}/internal/common/responses"
	"{{ cookiecutter.project_slug }}/internal/core/services"
	"{{ cookiecutter.project_slug }}/internal/core/vo"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(
	userService services.UserService,
) UserController {
	return UserController{
		userService: userService,
	}
}

func (u UserController) SignUp(c *gin.Context) {
	data := vo.UserVO{}
	responses.Ok(c, data)
}
