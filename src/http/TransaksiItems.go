package http

import (
	"github.com/gofiber/fiber/v2"
	"loundry/api/src/database"
	"loundry/api/src/helper"
	"loundry/api/src/models"
)

type TransaksiItems struct {
	ID          uint  `json:"id"`
	TransaksiID uint  `json:"transaksi_id"`
	ItemID      uint  `json:"item_id"`
	Item        Items `json:"item"`
	Qty         int   `json:"qty"`
	Total       int   `json:"harga"`
}

func TransaksiItemResponse(transaksiItem models.TransaksiItems) TransaksiItems {
	var item models.Items
	database.DBConn.First(&item, transaksiItem.IdItem)

	return TransaksiItems{
		ID:          transaksiItem.ID,
		TransaksiID: transaksiItem.TransaksiID,
		ItemID:      transaksiItem.IdItem,
		Qty:         transaksiItem.Qty,
		Total:       transaksiItem.Total,
		Item:        responseItem(item),
	}
}

func GetTransaksiitemsbyIdTransaction(ctx *fiber.Ctx) error {
	transactionId, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "transaction id not found",
		})
	}
	var transaksiItems []models.TransaksiItems
	if err := database.DBConn.Where("transaksi_id = ?", transactionId).Find(&transaksiItems).Error; err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "transaction id not found",
		})
	}
	var transaksiItemsResponse []TransaksiItems
	for _, transaksiItem := range transaksiItems {
		transaksiItemsResponse = append(transaksiItemsResponse, TransaksiItemResponse(transaksiItem))
	}
	return ctx.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   transaksiItemsResponse,
	})
}

func AddItemToTransaction(ctx *fiber.Ctx) error {
	transactionId, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "transaction id not found",
		})
	}
	userId := helper.GetUserIdFromToken(ctx)
	var transaksi models.Transaksi
	if err := database.DBConn.Where("id = ? AND user_id=?", transactionId, userId).Find(&transaksi).Error; err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "transaction not found",
		})
	}

	type BodyPost struct {
		ItemId int `json:"item_id"`
		Qty    int `json:"qty"`
	}
	var body BodyPost
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid body",
		})
	}
	var item models.Items
	userId = helper.GetUserIdFromToken(ctx)
	if err := findItemById(body.ItemId, userId, &item); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "item not found",
		})
	}

	var transaksiItem models.TransaksiItems
	var counter int64
	//check itransactionItemCOunt
	if err := database.DBConn.Model(&transaksiItem).Where("transaksi_id = ? AND id_item = ?", transactionId, body.ItemId).Count(&counter).Error; err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "transaction id not found",
		})
	}
	if counter > 0 {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "item already exist",
		})
	}

	newItem := models.TransaksiItems{
		TransaksiID: transaksi.ID,
		IdItem:      uint(body.ItemId),
		Total:       item.HargaBarang * body.Qty,
		Qty:         body.Qty,
	}
	if err := database.DBConn.Create(&newItem).Error; err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "transaction id not found",
		})
	}

	var totalTransaksi int
	var transaksiItems []models.TransaksiItems
	if err := database.DBConn.Where("transaksi_id = ?", transactionId).Find(&transaksiItems).Error; err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "transaction id not found",
		})
	}
	for _, transaksiItem := range transaksiItems {
		totalTransaksi += transaksiItem.Total
	}
	Updatedata := models.Transaksi{
		SubTotal: float64(totalTransaksi),
		Total:    float64(totalTransaksi) - transaksi.Diskon,
	}
	if err := database.DBConn.Model(&transaksi).Where("id = ?", transactionId).Updates(Updatedata).Error; err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "transaction id not found",
		})
	}
	return ctx.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   TransaksiItemResponse(newItem),
	})
}

func UpdateItemsTransaction(ctx *fiber.Ctx) error {
	transactionItemId, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "transaction item id not found",
		})
	}
	var transaksiItem models.TransaksiItems
	if err := database.DBConn.Where("id = ?", transactionItemId).Find(&transaksiItem).Error; err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "transaction item not found",
		})
	}
	type BodyPost struct {
		Qty int `json:"qty"`
	}

	var body BodyPost
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid body",
		})
	}
	var item models.Items
	if err := database.DBConn.Where("id = ?", transaksiItem.IdItem).Find(&item).Error; err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "item not found",
		})
	}
	newDataUpdate := models.TransaksiItems{
		Qty:   body.Qty,
		Total: body.Qty * item.HargaBarang,
	}
	if err := database.DBConn.Model(&transaksiItem).Updates(newDataUpdate).Error; err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to update transaction",
		})
	}
	var transaksi models.Transaksi
	if err := database.DBConn.Where("id = ?", transaksiItem.TransaksiID).Find(&transaksi).Error; err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "transaction not found",
		})
	}

	transactionItems := []models.TransaksiItems{}
	if err := database.DBConn.Where("transaksi_id = ?", transaksiItem.TransaksiID).Find(&transactionItems).Error; err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "transaction item not found",
		})
	}
	var total float64
	for _, item := range transactionItems {
		total += float64(item.Total)
	}
	subtotal := total - transaksi.Diskon
	dataUpdate := models.Transaksi{
		SubTotal: total,
		Total:    subtotal,
	}

	if err := database.DBConn.Model(&transaksi).Updates(dataUpdate).Error; err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to update transaction",
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   TransaksiItemResponse(transaksiItem),
	})
}

func DeleteTransactionItems(ctx *fiber.Ctx) error {
	transactionItemid, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "failed",
			"message": "uncompleted params",
		})
	}

	var transactionItem models.TransaksiItems
	if err := database.DBConn.Where("id=?", transactionItemid).Find(&transactionItem).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "failed",
			"message": "item not found",
		})
	}

	var transactionData models.Transaksi
	userId := helper.GetUserIdFromToken(ctx)
	if err := database.DBConn.Where("id=? AND user_id=?", transactionItem.TransaksiID, userId).Find(&transactionData).Error; err != nil {
		return ctx.Status(200).JSON(fiber.Map{
			"status":  "failed",
			"message": "forbidden",
		})
	}

	if err := database.DBConn.Delete(&transactionItem).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "failed",
			"message": "server error",
		})
	}

	var transactionItems []models.TransaksiItems
	if err := database.DBConn.Where("transaksi_id=?", transactionData.ID).Find(&transactionItems).Error; err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"status":  "failed",
			"message": "server error",
		})
	}
	var total int
	for _, itemTransaksi := range transactionItems {
		total += itemTransaksi.Total
	}
	dataUpdate := models.Transaksi{
		SubTotal: float64(total),
		Total:    float64(total) - transactionData.Diskon,
	}

	if err := database.DBConn.Where("id", transactionData.ID).Updates(dataUpdate).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "failed",
			"message": "unable to update transaction data",
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "item has been deleted",
	})
}
