package http

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"loundry/api/src/database"
	"loundry/api/src/helper"
	"loundry/api/src/models"
	"time"
)

func AuthCheck(ctx *fiber.Ctx) error {
	type AuthInput struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type Response struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Nama     string `json:"nama"`
		Token    string `json:"token"`
	}
	user := []models.Users{}
	dataUser := new(AuthInput)
	var responseData Response

	if err := ctx.BodyParser(dataUser); err != nil {
		fmt.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
		})
	}

	if err := database.DBConn.Where("username = ?", dataUser.Username).First(&user).Error; err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.ErrUnauthorized,
			"message": "Oops! Something went wrong",
		})
	}
	pass := dataUser.Password
	err := bcrypt.CompareHashAndPassword([]byte(user[0].Password), []byte(pass))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.ErrUnauthorized,
			"message": "Oops! Something went wrong",
		})
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = dataUser.Username
	claims["user_id"] = user[0].ID
	claims["nama"] = user[0].Nama
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	tokenString, _ := token.SignedString([]byte(helper.ReadEnv("JWT_SECRET")))
	responseData = Response{
		ID:       user[0].ID,
		Username: user[0].Username,
		Nama:     user[0].Nama,
		Token:    tokenString,
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Login Success",
		"data":    responseData,
	})
}
