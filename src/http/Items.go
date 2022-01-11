package http

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"loundry/api/src/database"
	"loundry/api/src/helper"
	"loundry/api/src/models"
)

type Items struct {
	ID         uint   `json:"id"`
	NamaBarang string `json:"nama_barang"`
	Harga      int    `json:"harga"`
}

func responseItem(item models.Items) Items {
	return Items{
		ID:         item.ID,
		NamaBarang: item.NamaBarang,
		Harga:      item.HargaBarang,
	}
}

func findItemById(id int, userId uint, item *models.Items) error {
	database.DBConn.Where("id = ? AND user_id = ?", id, userId).Find(&item)
	if item.ID == 0 {
		return errors.New("Item not found")
	}

	return nil
}

func FindItemById(id uint, userId uint, item *models.Items) error {
	database.DBConn.Where("id = ? AND user_id = ?", id, userId).Find(&item)
	if item.ID == 0 {
		return errors.New("Item not found")
	}

	return nil
}

func CreateNewItems(ctx *fiber.Ctx) error {
	type RequestBody struct {
		NamaBarang string `json:"nama_barang"`
		Harga      int    `json:"harga"`
	}

	var body RequestBody
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}
	userId := helper.GetUserIdFromToken(ctx)

	var count int64
	database.DBConn.Model(&models.Items{}).Where("nama_barang = ?", body.NamaBarang).Count(&count)
	if count > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Nama barang sudah ada",
		})
	}
	newItem := models.Items{
		NamaBarang:  body.NamaBarang,
		HargaBarang: body.Harga,
		UserId:      userId,
	}

	database.DBConn.Create(&newItem)
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Item berhasil ditambahkan",
		"data":    responseItem(newItem),
	})
}

func GetAllItems(ctx *fiber.Ctx) error {
	var items []models.Items
	database.DBConn.Find(&items)

	var response []Items
	for _, item := range items {
		response = append(response, responseItem(item))
	}
	return ctx.JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}

func GetItemById(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	var item models.Items
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}
	userId := helper.GetUserIdFromToken(ctx)
	if err := findItemById(id, userId, &item); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	return ctx.JSON(fiber.Map{
		"status": "success",
		"data":   responseItem(item),
	})
}

func UpdateItemById(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}
	var item models.Items
	userId := helper.GetUserIdFromToken(ctx)
	if err := findItemById(id, userId, &item); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	type RequestBody struct {
		Harga int `json:"harga"`
	}

	var body RequestBody
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	database.DBConn.Model(&item).Updates(models.Items{
		HargaBarang: body.Harga,
	})

	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": "Item berhasil diupdate",
		"data":    responseItem(item),
	})
}

func DeleteItemById(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}
	var item models.Items
	userId := helper.GetUserIdFromToken(ctx)
	if err := findItemById(id, userId, &item); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	database.DBConn.Delete(&item)
	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": "Item berhasil dihapus",
	})
}
