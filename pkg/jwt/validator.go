package jwt

import (
	"errors"
	"giiku-camp/internal/domain/entity"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// トークン文字列が有効化検証する
func VerifyToken(tokenStr string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET is not set")
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		// HS256 で来ているかチェック
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, entity.ErrTokenInValid
	}

	return token, nil
}

// リクエストヘッダーからトークン文字列を取得し、検証する
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

// リクエストヘッダーからトークン文字列を取得し、ユーザーIDを取得する
func ExtractUserIDFromContext(c *gin.Context) (string, error) {
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
		return "", entity.ErrFailedToParseClaims
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", entity.ErrFailedToParseClaims
	}

	return userID, nil
}

// トークンからユーザーIDを取得する
func ExtractUserIDFromToken(token *jwt.Token) (string, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", entity.ErrFailedToParseClaims
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", entity.ErrFailedToParseClaims
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
