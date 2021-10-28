package notification

import (
	"fmt"
	"github.com/Mikkelhost/Gophers-Honey/pkg/config"
	"github.com/Mikkelhost/Gophers-Honey/pkg/model"
	"github.com/rs/zerolog/log"
	"net/smtp"
	"strconv"
)

var smtpServer model.SmtpServer

func isSmtpServerConfigured() bool {
	return false
}

func ConfigureSmtpServer(port uint16, username, password, mailserver string) {
	var configSmtpServer model.SmtpServer
	configSmtpServer.Username = username
	configSmtpServer.Password = password
	configSmtpServer.SmtpHost = mailserver
	configSmtpServer.SmtpPort = port

	config.SetSmtpServer(configSmtpServer)
}

func getSmtpServer() {
	var temp model.SmtpServer

	smtpServer = temp
}

func constructMessage(alert model.Log) []byte {
	// TODO: Write prefix message
	prefix := ""
	message := fmt.Sprintf("%s! Device with Device ID: %d, has on %s raised alert based on the message %s",
		alert.Level, alert.DeviceID, alert.TimeStamp.String(), alert.Message)

	byteMessage := []byte(prefix + message)
	return byteMessage
}

func SendEmailNotification(alert model.Log, to []string) error {
	message := constructMessage(alert)
	from := smtpServer.Username

	stringPort := strconv.Itoa(int(smtpServer.SmtpPort))
	auth := smtp.PlainAuth("", smtpServer.Username, smtpServer.Password, smtpServer.SmtpHost)

	err := smtp.SendMail(smtpServer.SmtpHost+":"+stringPort, auth, from, to, message)
	if err != nil {
		log.Logger.Warn().Msgf("Error sending email: %s", err)
		return err
	}

	return nil
}
