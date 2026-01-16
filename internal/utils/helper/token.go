package helper

import (
	"ngevent/internal/model"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(user *model.Users) (string, string, error) {
	// Access Token
	token1 := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"user": user.Email,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	accessToken, err := token1.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", "", nil
	}

	// Refresh Token
	token2 := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"user": user.Email,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	refreshToken, err := token2.SignedString([]byte(os.Getenv("REFRESH_SECRET_KEY")))
	if err != nil {
		return "", "", nil
	}

	return accessToken, refreshToken, nil
}
