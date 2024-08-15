package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	"{{ cookiecutter.project_slug }}/configs"
	"{{ cookiecutter.project_slug }}/internal/core/validation"
	"{{ cookiecutter.project_slug }}/internal/helpers/auth"
	"{{ cookiecutter.project_slug }}/internal/helpers/responses"
)

func JwtAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader("Authorization")

		if authorizationHeader == "" {
			responses.AbortWithAPIError(ctx, validation.ErrAccessTokenInvalid)
			return
		}

		fields := strings.Fields(authorizationHeader)
		tokenString := fields[1]

		claims, err := auth.VerifyJwtToken(configs.Env.JwtSecret, tokenString)
		if err != nil {
			responses.AbortWithAPIError(ctx, validation.ErrAccessTokenInvalid)
			return
		}

		ctx.Set("userID", claims.UserID.String())

		ctx.Next()
	}
}
