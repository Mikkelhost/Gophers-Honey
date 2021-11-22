package notification

import (
	"fmt"
	"github.com/Mikkelhost/Gophers-Honey/pkg/config"
	"github.com/Mikkelhost/Gophers-Honey/pkg/database"
	"github.com/Mikkelhost/Gophers-Honey/pkg/model"
	"github.com/rs/zerolog/log"
	"net/smtp"
	"strconv"
)

// ConfigureSmtpServer sets SMTP server and user configuration and writes
// changes to the configuration file.
func ConfigureSmtpServer(configSmtpServer model.SmtpServer) error {

	config.Conf.SmtpServer = configSmtpServer // Set configuration in memory.
	err := config.WriteConf()                 // Write configuration to config file.
	if err != nil {
		log.Logger.Warn().Msgf("Error writing SMTP server configuration to config file")
		return err
	}

	return nil
}

// constructMessage uses the alert to construct a message.
func constructMessage(alert model.Log) []byte {
	// TODO: Write prefix message
	prefix := ""
	message := fmt.Sprintf("Device with Device ID: %d, has on %s raised alert based on the message %s",
		alert.DeviceID, alert.TimeStamp.String(), alert.Message)

	byteMessage := []byte(prefix + message)
	return byteMessage
}

// SendEmailNotification handles the construction of email messages as
// well as sending the constructed messages.
func SendEmailNotification(alert model.Log, to []string) error {
	message := constructMessage(alert)
	smtpServer := config.Conf.SmtpServer
	from := smtpServer.Username

	stringPort := strconv.Itoa(int(smtpServer.SmtpPort))
	auth := smtp.PlainAuth("", smtpServer.Username, smtpServer.Password, smtpServer.SmtpHost)

	err := smtp.SendMail(smtpServer.SmtpHost+":"+stringPort, auth, from, to, message)
	if err != nil {
		log.Logger.Warn().Msgf("Error sending email: %s", err)
		return err
	}
	log.Logger.Info().Msgf("Email sent")

	return nil
}

// NotifyAll fetches the email addresses of users with notifications
// enabled and sends a mail with the alert.
func NotifyAll(alert model.Log) error {
	var emails []string

	users, err := database.GetAllUsers()
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.NotificationsEnabled {
			emails = append(emails, user.Email)
		}
	}

	err = SendEmailNotification(alert, emails)
	if err != nil {
		return err
	}

	return nil
}
