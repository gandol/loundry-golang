package helper

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func GetUserIdFromToken(ctx *fiber.Ctx) uint {
	var user *jwt.Token
	l := ctx.Locals("user")
	if l == nil {
		return 0
	}
	user = l.(*jwt.Token)
	id := uint(((user.Claims.(jwt.MapClaims)["user_id"]).(float64)))
	return id
}
