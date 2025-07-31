package middleware

import (
	"fmt"
	"strings"

	"ipincamp/srikandi-sehat/config"
	"ipincamp/srikandi-sehat/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Missing authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid token format, 'Bearer ' prefix missing")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Get("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid token claims")
	}

	userUUID, ok := claims["usr"].(string)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid user identifier in token")
	}

	c.Locals("usr", userUUID)
	return c.Next()
}
