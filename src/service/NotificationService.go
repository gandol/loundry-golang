package service

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"loundry/api/src/ENUM"
	"loundry/api/src/database"
	"loundry/api/src/helper"
	"loundry/api/src/models"
)

func CreateNotification(transaksi models.Transaksi, notificationTeks string) error {
	var customerData models.Customer
	if err := database.DBConn.Where("id = ?", transaksi.CustomerId).First(&customerData).Error; err != nil {
		return err
	}
	//get setting data
	var settingData models.Settings
	if err := database.DBConn.Where("name = ?", ENUM.WHATSAPP_NUMBER).First(&settingData).Error; err != nil {
		return err
	}

	newData := models.Notifications{
		TransactionId:    transaksi.ID,
		NotificationTeks: notificationTeks,
		HasSent:          false,
		PhoneNumber:      customerData.PhoneNumber,
	}
	if err := SendNotifToBroker(notificationTeks, settingData.Value, customerData.PhoneNumber); err != nil {
		return err
	}

	if err := database.DBConn.Create(&newData).Error; err != nil {
		return err
	}

	return nil
}

type PesanKirim struct {
	User        string
	Message     string
	PhoneClient string
}

func SendNotifToBroker(message string, User string, PhoneClient string) error {
	amqpServer := helper.ReadEnv("AMQP_SERVER")
	connectRabitMq, err := amqp.Dial(amqpServer)
	if err != nil {
		return err
	}
	defer connectRabitMq.Close()
	amqpChannel, err := connectRabitMq.Channel()
	if err != nil {
		return err
	}
	defer amqpChannel.Close()

	_, err = amqpChannel.QueueDeclare(
		User,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	PesanKirim := PesanKirim{
		User:        User,
		Message:     message,
		PhoneClient: PhoneClient,
	}
	pesanString, err := json.Marshal(PesanKirim)
	pesan := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(string(pesanString)),
	}
	err = amqpChannel.Publish(
		"",
		User,
		false,
		false,
		pesan,
	)
	if err != nil {
		return err
	}
	return nil
}
