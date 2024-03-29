package middleware

import (
	"InternBCC/entity"
	"InternBCC/sdk"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
	"time"
)

func Auth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		sdk.FailOrError(c, http.StatusUnauthorized, "Failed to get token", err)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//cek expired
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "fail",
				"message": "Token expired",
			})
			return
		}

		//cari user
		//var user entity.UserRegister
		//database.DB.First(&user, claims["sub"])
		//if user.ID == 0 {
		//	c.JSON(400, gin.H{
		//		"error": "error 3",
		//	})
		//	return
		//}
		c.Set("user", claims["sub"])
		c.Next()

	} else {
		c.JSON(401, gin.H{
			"status":  "fail",
			"message": "Token is not valid",
		})
		return
	}
}
func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authorization := c.Request.Header.Get("Authorization")
		if !strings.HasPrefix(authorization, "Bearer ") {
			c.Abort()
			msg := "wrong header value"
			sdk.FailOrError(c, http.StatusForbidden, msg, errors.New(msg))
			return
		}
		tokenJwt := authorization[7:]
		claims := entity.UserClaims{}
		jwtKey := os.Getenv("SECRET_KEY")
		if _, err := sdk.DecodeToken(tokenJwt, &claims, jwtKey); err != nil {
			c.Abort()
			sdk.FailOrError(c, http.StatusUnauthorized, "unauthorized", err)
			return
		}
		c.Set("user", claims)
	}
}
