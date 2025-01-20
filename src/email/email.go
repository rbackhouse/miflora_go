package email

import (
	"fmt"
	"net/smtp"

	logger "github.com/sirupsen/logrus"

	"potpie.org/miflora/src/config"
)

type email struct {
	toAddresses []string
	config      config.EmailConfig
}

func sendEmail(config config.EmailConfig, message string, toAddresses []string) {
	auth := smtp.PlainAuth("", config.FromAddress, config.Password, config.SmtpHost)
	host := fmt.Sprintf("%s:%d", config.SmtpHost, config.SmtpPort)

	err := smtp.SendMail(host, auth, config.FromAddress, toAddresses, []byte(message))
	if err != nil {
		logger.Warnf("Failed to send email: %s", err.Error())
	}
}
