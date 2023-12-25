package jwt

import (
	"net/http"
	"strings"
	"utils/conf"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	UserID   string
	Username string
	jwt.RegisteredClaims
}

func CreateJWT(id string, name string) (string, error) {
	claims := UserClaims{
		UserID:   id,
		Username: name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 0, 7)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    conf.Conf().JWT_ISSUER,
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenClaims.SignedString([]byte(conf.Conf().JWT_SECRET))
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
