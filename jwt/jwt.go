package jwt

import (
	"chat_room/conf"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
)

type UserClaims struct {
	UserID   string
	Username string
	jwt.RegisteredClaims
}

func ParseJWT(token string) (*UserClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.Conf().JWT_SECRET), nil
	})
	if err == nil && jwtToken != nil {
		if claim, ok := jwtToken.Claims.(*UserClaims); ok && jwtToken.Valid {
			return claim, nil
		}
	}
	return nil, err
}
func AuthenticateJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization := ctx.GetHeader("Authorization")
		if !strings.HasPrefix(authorization, "Bearer ") {
			ctx.Abort()
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorization header",
			})
			return
		}
		token := authorization[7:]
		claim, err := ParseJWT(token)
		ctx.Set("jwt", claim)
		if err != nil {
			ctx.Abort()
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.Next()
	}
}
