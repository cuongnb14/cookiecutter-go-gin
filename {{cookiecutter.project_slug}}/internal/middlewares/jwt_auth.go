package middlewares

import (
	"github.com/gin-gonic/gin"
	"{{ cookiecutter.project_slug }}/configs"
	"{{ cookiecutter.project_slug }}/internal/utils/responses"
	jwt_utils "{{ cookiecutter.project_slug }}/pkg/jwt"
)

func JwtAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader("Authorization")

		if authorizationHeader == "" {
			responses.AbortWithAPIError(ctx, responses.ErrAccessTokenInvalid)
			return
		}

		fields := strings.Fields(authorizationHeader)
		tokenString := fields[1]

		claims, err := auth.VerifyJwtToken(configs.Env.JwtSecret, tokenString)
		if err != nil {
			responses.AbortWithAPIError(ctx, responses.ErrAccessTokenInvalid)
			return
		}
		ctx.Set("userID", claims.UserID.String())
		ctx.Next()
	}
}
