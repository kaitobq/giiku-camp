package jwt

import (
	"errors"
	"fmt"
	"giiku-camp/internal/domain/entity"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateAccessToken(userID string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET is not set")
	}

	accessTokenLifeStr := os.Getenv("ACCESS_TOKEN_LIFE_SPAN")
	if accessTokenLifeStr == "" {
		return "", errors.New("ACCESS_TOKEN_LIFE_SPAN is not set")
	}
	accessTokenLife, err := strconv.Atoi(accessTokenLifeStr)
	if err != nil {
		return "", fmt.Errorf("invalid ACCESS_TOKEN_LIFE_SPAN: %w", err)
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * time.Duration(accessTokenLife)).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func GenerateRefreshToken(userID string, tokenVersion int) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET is not set")
	}

	refreshTokenLifeStr := os.Getenv("REFRESH_TOKEN_LIFE_SPAN")
	if refreshTokenLifeStr == "" {
		return "", errors.New("REFRESH_TOKEN_LIFE_SPAN is not set")
	}
	refreshTokenLife, err := strconv.Atoi(refreshTokenLifeStr)
	if err != nil {
		return "", fmt.Errorf("invalid REFRESH_TOKEN_LIFE_SPAN: %w", err)
	}

	claims := jwt.MapClaims{
		"user_id":       userID,
		"token_version": tokenVersion,
		"exp":           time.Now().Add(time.Hour * time.Duration(refreshTokenLife)).Unix(),
		"iat":           time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func RefreshTokens(user entity.User, tokenStr string) (string, string, error) {
	token, err := VerifyToken(tokenStr)
	if err != nil {
		return "", "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", entity.ErrFailedToParseClaims
	}
	tokenVersion, _ := claims["token_version"].(float64)
	if !ok {
		return "", "", entity.ErrFailedToParseClaims
	}

	if int(tokenVersion) != user.TokenVersion {
		return "", "", entity.ErrTokenVersionMismatch
	}

	accessToken, err := GenerateAccessToken(user.ID)
	if err != nil {
		return "", "", err
	}
	user.IncrementTokenVersion()
	refreshToken, err := GenerateRefreshToken(user.ID, user.TokenVersion)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
