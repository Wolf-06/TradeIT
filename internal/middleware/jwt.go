package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var Key = []byte(os.Getenv("jwtkey"))

func CreateToken(userid uint64) string {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userid,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenstring, err := token.SignedString(Key)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Println(tokenstring)
	return tokenstring
}

func VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization Token"})
			return
		}
		Token, err := jwt.Parse(header, func(Token *jwt.Token) (interface{}, error) {
			if _, ok := Token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid token signing method")
			}
			return Key, nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			log.Fatalln("Invalid Token ", err)
			return
		}

		if claim, ok := Token.Claims.(jwt.MapClaims); ok && Token.Valid {
			c.Set("userid", claim["sub"])
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
		}
	}
}
