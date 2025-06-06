package auth

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Pleiades-IUST/backend/utils/config"
	"github.com/Pleiades-IUST/backend/utils/ginutil"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthenticateUser(ctx *gin.Context) {
	hmacSecret := []byte(config.GetSecretKey())

	authHeader := ctx.Request.Header.Get("Authorization")
	bearerToken := ""
	if strings.HasPrefix(authHeader, "Bearer ") {
		bearerToken = strings.TrimPrefix(authHeader, "Bearer ")
	} else {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userIDStr, err := parseLoginToken(bearerToken, hmacSecret)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return
	}

	ctx.Set(ginutil.UserIDKey, userID)
}

func parseLoginToken(bearerToken string, hmacSecret []byte) (string, error) {
	token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})
	if err != nil {
		return "", err
	}

	return token.Claims.GetSubject()
}
