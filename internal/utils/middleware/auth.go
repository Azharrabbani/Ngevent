package middleware

import (
	"fmt"
	"ngevent/internal/dto"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Sub  string `json:"sub"`
	User string `json:"user"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Retrive cookie
		tokenString := c.Cookies("ngevent-cookie")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.Error(
				fiber.StatusUnauthorized,
				"error",
				"unauthorized",
				"missing jwt cookie",
			))
		}

		// Verifiy the cookie
		claims, err := VerifyToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.Error(
				fiber.StatusUnauthorized,
				"error",
				"unauthorized",
				"invalid token",
			))
		}

		c.Locals("user_id", claims.Sub)
		c.Locals("role", claims.Role)
		c.Locals("exp", claims.ExpiresAt)

		return c.Next()
	}
}

func VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
