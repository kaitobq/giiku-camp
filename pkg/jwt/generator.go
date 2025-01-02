package jwt

import (
	"errors"
	"fmt"
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

func GenerateRefreshToken(userID string) (string, error) {
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
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * time.Duration(refreshTokenLife)).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

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

func RefreshTokens(oldRefreshToken string) (string, string, error) {
	// 古いリフレッシュトークンを検証
	token, err := VerifyToken(oldRefreshToken)
	if err != nil {
		return "", "", fmt.Errorf("refresh token invalid: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", "", errors.New("invalid refresh token")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", "", errors.New("missing user_id in token claims")
	}

	// 新しいアクセストークンとリフレッシュトークンを発行
	newAccessToken, err := GenerateAccessToken(userID)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}
	newRefreshToken, err := GenerateRefreshToken(userID)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return newAccessToken, newRefreshToken, nil
}
