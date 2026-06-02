package helper

import (
	"errors"
	"fmt"
	"ngevent/internal/model"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateAccessToken(user *model.Users) (string, error) {
	now := time.Now()

	claims := jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"iat":  now.Unix(),
		"nbf":  now.Unix(),
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("SECRET_KEY")))
}

func GenerateRefreshToken(userID string, expire time.Time) (string, string, error) {
	now := time.Now()

	jti := uuid.NewString() // unique token id

	claims := jwt.MapClaims{
		"sub": userID,
		"jti": jti,
		"iat": now.Unix(),
		"nbf": now.Unix(),
		"exp": expire.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_REFRESH_KEY")))
	if err != nil {
		return "", "", err
	}

	return tokenString, jti, nil
}

func ValidateRefreshToken(refreshToken string) (string, string, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return []byte(os.Getenv("SECRET_REFRESH_KEY")), nil
	})

	if err != nil || !token.Valid {
		return "", "", errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("invalid token claims")
	}

	jtiRaw, ok := claims["jti"]
	if !ok {
		return "", "", errors.New("missing jti")
	}

	jti, ok := jtiRaw.(string)
	if !ok {
		return "", "", errors.New("invalid jti format")
	}

	userID := claims["sub"].(string)

	return userID, jti, nil
}

func ClearAuthCookies(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     "ngevent_cookie",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
	})
}
