package jwt

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func VerifyToken(tokenStr string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET is not set")
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		// HS256 で来ているかチェック
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("token is invalid or expired")
	}

	return token, nil
}

func VerifyReqHeaderToken(c *gin.Context) (bool, error) {
	tokenStr, err := getTokenStringFromRequestHeader(c)
	if err != nil {
		return false, err
	}

	_, err = VerifyToken(tokenStr)
	if err != nil {
		return false, err
	}

	return true, nil
}

func ExtractUserIDFromToken(c *gin.Context) (string, error) {
	tokenStr, err := getTokenStringFromRequestHeader(c)
	if err != nil {
		return "", err
	}
	token, err := VerifyToken(tokenStr)
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("failed to parse claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("failed to parse user_id")
	}

	return userID, nil
}

func getTokenStringFromRequestHeader(c *gin.Context) (string, error) {
	bearToken := c.Request.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1], nil
	}

	return "", errors.New("no token found")
}
