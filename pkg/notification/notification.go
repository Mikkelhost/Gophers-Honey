package notification

import (
	"fmt"
	"github.com/Mikkelhost/Gophers-Honey/pkg/config"
	"github.com/Mikkelhost/Gophers-Honey/pkg/database"
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"github.com/Mikkelhost/Gophers-Honey/pkg/model"
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
	// TODO: Add additional info to message.
	prefix := ""
	message := fmt.Sprintf( "Subject: Test email from gopher\r\n" +
		"\r\n" +
		"Device with Device ID: %d, has on %s raised alert based on the message %s",
		alert.DeviceID, alert.LogTimeStamp.String(), alert.Message)

	byteMessage := []byte(prefix + message)
	return byteMessage
}

// SendEmailNotification handles the construction of email messages as
// well as sending the constructed messages.
// TODO: need to present alert message in a nice way.
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

func SendTestEmail(to []string) error {
	message := "To: "+to[0]+"\r\n" +
		"Subject: Test email from gopher\r\n" +
		"\r\n" +
		"This is a test email sent from your gophers-honey setup. If you have received this email you have set up the email configuration correctly"
	smtpServer := config.Conf.SmtpServer
	from := smtpServer.Username

	stringPort := strconv.Itoa(int(smtpServer.SmtpPort))
	auth := smtp.PlainAuth("", smtpServer.Username, smtpServer.Password, smtpServer.SmtpHost)

	err := smtp.SendMail(smtpServer.SmtpHost+":"+stringPort, auth, from, to, []byte(message))
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
