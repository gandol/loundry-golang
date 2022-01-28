package http

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io/ioutil"
	"loundry/api/src/ENUM"
	"loundry/api/src/database"
	"loundry/api/src/helper"
	"loundry/api/src/models"
	"net/http"
	"strings"
)

type Settings struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"descriptsion"`
	Value       string `json:"value"`
}

func responseSetting(setting models.Settings) Settings {
	return Settings{
		ID:          setting.ID,
		Name:        setting.Name,
		Description: setting.Description,
		Value:       setting.Value,
	}
}

func CreateNotificationSettings(ctx *fiber.Ctx) error {
	type RequestUser struct {
		Value string `json:"value"`
	}
	var body RequestUser
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	userId := helper.GetUserIdFromToken(ctx)
	var count int64
	if err := database.DBConn.Model(&models.Settings{}).Where("name = ? AND user_id=? ", ENUM.WHATSAPP_NUMBER, userId).Count(&count).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if count > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Setting already exist",
		})
	}
	setting := models.Settings{
		Value:       body.Value,
		Description: "",
		Name:        ENUM.WHATSAPP_NUMBER,
		UserID:      userId,
	}
	if err := database.DBConn.Create(&setting).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	payload := strings.NewReader(`{"phoneNumber":"` + string(body.Value) + `",
    "isActive":0,
    "qrString":""}`)
	resp, error := http.Post("http://localhost:8001/users", "application/json", payload)
	if error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": error.Error(),
		})
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(respBody))

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Setting created",
		"data":    responseSetting(setting),
	})
}

func UpdateNotificationsSetting(ctx *fiber.Ctx) error {
	type RequestUser struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
	var body RequestUser
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	userId := helper.GetUserIdFromToken(ctx)
	var setting models.Settings
	if err := database.DBConn.Where("name = ? AND user_id=? ", body.Name, userId).First(&setting).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	setting.Value = body.Value
	if err := database.DBConn.Save(&setting).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Setting updated",
		"data":    responseSetting(setting),
	})
}
func GetNotificationNotificationQris(ctx *fiber.Ctx) error {
	type dataQr struct {
		QrString string `json:"qrString"`
	}
	userId := helper.GetUserIdFromToken(ctx)
	var setting models.Settings
	if err := database.DBConn.Where("name = ? AND user_id=? ", ENUM.WHATSAPP_NUMBER, userId).First(&setting).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	resp, errs := http.Get("http://localhost:8001/users/qr/" + setting.Value)
	if errs != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error",
		})
	}
	defer resp.Body.Close()
	var data dataQr
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Setting found",
		"data":    data,
	})
}
