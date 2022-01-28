package http

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"loundry/api/src/database"
	"loundry/api/src/helper"
	"loundry/api/src/models"
	"loundry/api/src/service"
)

const (
	Pending           = "Pending"
	ProsesCuci        = "Proses Cuci"
	ProsesPengeringan = "Proses Pengeringan"
	ProsesSetrika     = "Proses Setrika"
	Selesai           = "Selesai"
)

type Transaksi struct {
	ID             uint             `json:"id"`
	CustomerId     uint             `json:"customer"`
	SubTotal       float64          `json:"subtotal"`
	Diskon         float64          `json:"diskon"`
	Total          float64          `json:"total"`
	Status         string           `json:"status"`
	TransaksiItems []TransaksiItems `json:"items"`
}

//
func transaksiResponse(transaksi models.Transaksi) Transaksi {
	transaksiItems := []models.TransaksiItems{}
	database.DBConn.Where("transaksi_id = ?", transaksi.ID).Find(&transaksiItems)
	var dataItems = []TransaksiItems{}
	for _, item := range transaksiItems {
		var itemData models.Items
		database.DBConn.Where("id = ?", item.IdItem).First(&itemData)
		dataItems = append(dataItems, TransaksiItems{
			ID:          item.ID,
			TransaksiID: item.TransaksiID,
			ItemID:      item.IdItem,
			Qty:         item.Qty,
			Total:       item.Total,
			Item:        responseItem(itemData),
		})
	}
	return Transaksi{
		ID:             transaksi.ID,
		CustomerId:     transaksi.CustomerId,
		SubTotal:       transaksi.SubTotal,
		Diskon:         transaksi.Diskon,
		Total:          transaksi.Total,
		Status:         transaksi.Status,
		TransaksiItems: dataItems,
	}
}

//
func findTransaksiById(id int, userId uint, transaksi *models.Transaksi) error {
	database.DBConn.Where("id = ? AND user_id = ?", id, userId).Find(&transaksi)
	if transaksi.ID == 0 {
		return errors.New("Transaksi not found")
	}
	return nil
}

func GetAllTransaction(ctx *fiber.Ctx) error {
	var transaksi []models.Transaksi
	userId := helper.GetUserIdFromToken(ctx)
	if err := database.DBConn.Where("user_id=?", userId).Find(&transaksi).Error; err != nil {
		fmt.Println(err.Error)
		return ctx.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Internal Server Error",
		})
	}
	var dataTransaksi []Transaksi
	for _, item := range transaksi {
		dataTransaksi = append(dataTransaksi, transaksiResponse(item))
	}
	return ctx.JSON(fiber.Map{
		"status": "success",
		"data":   dataTransaksi,
	})
}

func GetTransactionById(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	var transaksi models.Transaksi
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}
	userId := helper.GetUserIdFromToken(ctx)
	if err := findTransaksiById(id, userId, &transaksi); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	return ctx.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   transaksiResponse(transaksi),
	})
}

func CreateTransaction(ctx *fiber.Ctx) error {
	type ItemsBody struct {
		ItemID int `json:"id"`
		Qty    int `json:"qty"`
	}
	type TransaksiParam struct {
		CustomerId uint        `json:"customer_id"`
		Items      []ItemsBody `json:"items"`
	}

	var transaksiParam TransaksiParam
	if err := ctx.BodyParser(&transaksiParam); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	userId := helper.GetUserIdFromToken(ctx)
	transaksi := models.Transaksi{
		CustomerId: transaksiParam.CustomerId,
		UserId:     userId,
	}
	if err := database.DBConn.Create(&transaksi).Error; err != nil {
		ctx.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Internal Server Error",
		})
	}
	var total int

	for _, item := range transaksiParam.Items {
		fmt.Println(item.ItemID)
		var itemModel models.Items
		if err := FindItemById(uint(item.ItemID), userId, &itemModel); err != nil {
			ctx.Status(404).JSON(err.Error())
		}

		transaksiItem := models.TransaksiItems{
			TransaksiID: transaksi.ID,
			IdItem:      uint(item.ItemID),
			Qty:         item.Qty,
			Total:       item.Qty * itemModel.HargaBarang,
		}
		if err := database.DBConn.Create(&transaksiItem).Error; err != nil {
			ctx.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Internal Server Error",
			})
		}
		total += transaksiItem.Total
	}
	transaksi.SubTotal = float64(total)
	transaksi.Total = float64(total)
	transaksi.Diskon = 0
	transaksi.Status = Pending
	if err := database.DBConn.Save(&transaksi).Error; err != nil {
		ctx.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Internal Server Error",
		})
	}

	return ctx.Status(201).JSON(fiber.Map{
		"status": "success",
		"data":   transaksiResponse(transaksi),
	})
}

func UpdateStatustransaksi(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "failde",
			"message": "Bad request",
		})
	}
	type DataUpdate struct {
		Status string `json:"status"`
	}

	var dataUpdate DataUpdate
	if err := ctx.BodyParser(&dataUpdate); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "failed",
			"message": "unclompleted data",
		})
	}
	var detailTransaksi models.Transaksi
	if err := database.DBConn.Where("id=?", id).Find(&detailTransaksi).Error; err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Internal Server Error",
		})
	}

	var dataTransaksi models.Transaksi
	switch dataUpdate.Status {
	case Pending:
		dataTransaksi.Status = Pending
		break
	case ProsesCuci:
		dataTransaksi.Status = ProsesCuci
		break
	case ProsesPengeringan:
		dataTransaksi.Status = ProsesPengeringan
		break
	case ProsesSetrika:
		dataTransaksi.Status = ProsesSetrika
		break
	case Selesai:
		dataTransaksi.Status = Selesai

		if err := service.CreateNotification(detailTransaksi, "Baju anda sudah selsai di loundry"); err != nil {
			return ctx.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Internal Server Error",
			})
		}
		break
	default:
		return ctx.Status(403).JSON(fiber.Map{
			"status":  "failed",
			"message": "error on the data",
		})
	}

	if err := database.DBConn.Where("id=?", id).Updates(&dataTransaksi).Error; err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"status":  "failed",
			"message": "unable to update data",
		})
	}
	return ctx.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Operation success",
	})
}
