package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// TODO: change into a certificate file bytes
var jwtKey = []byte("supersecretkey")

const TokenExpiration = 1 * time.Hour

func GenerateJWT(userId uint) (tokenString string, err error) {
	expirationTime := time.Now().Add(TokenExpiration)
	claims := &jwt.StandardClaims{
		Subject:   strconv.FormatUint(uint64(userId), 10),
		ExpiresAt: expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	return
}

func ParseToken(signedToken string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&jwt.StandardClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return nil, err
	}

	// Parse Claims interface into StandardClaims
	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		err = errors.New("could not parse claims")
		return nil, err
	}

	err = claims.Valid()
	if err != nil {
		return nil, err
	}
	return claims, nil

}

func VerifyID(c *gin.Context, id int) error {
	tokenSub := c.GetString("jwt_sub")
	tokenId, _ := strconv.Atoi(tokenSub)
	if id != tokenId {
		return errors.New("unauthorized id")
	}
	return nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"error":   "Unauthorized",
					"message": "No Access token provided.",
				},
			)
			return
		}
		splitToken := strings.Split(tokenString, "Bearer")
		if len(splitToken) != 2 {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"error":   "Unauthorized",
					"message": "Invalid bearer token format.",
				},
			)
			return
		}
		tokenString = strings.TrimSpace(splitToken[1])

		claims, err := ParseToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"error":   "Unauthorized",
					"message": err.Error(),
				},
			)
			return
		}
		// Set this key for scope authorizat
		c.Set("jwt_sub", claims.Subject)
		c.Next()
	}
}
