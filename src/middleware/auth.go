package middleware

import (
	_ "github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"loundry/api/src/helper"
)

func ProtectedArea() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(helper.ReadEnv("JWT_SECRET")),
		ErrorHandler: jwtError,
	})
}
func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})

	} else {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
	}
}

func RestrictedArea() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return ctx.Next()
	}
}
