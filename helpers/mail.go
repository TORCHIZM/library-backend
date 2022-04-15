package helpers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/firmanJS/fiber-with-mongo/config"

	mail "github.com/xhit/go-simple-mail/v2"
)

func SendMail(receiver string, code int) {
	server := mail.NewSMTPClient()
	server.Host = config.Config("MAIL_HOST")
	server.Port, _ = strconv.Atoi(config.Config("MAIL_PORT"))
	server.Username = config.Config("MAIL_ADDR")
	server.Password = config.Config("MAIL_PASS")

	server.Encryption = mail.EncryptionTLS

	smtpClient, err := server.Connect()
	if err != nil {
		log.Fatal(err)
	}

	email := mail.NewMSG()
	email.SetFrom("Kütüphane Uygulaması <dpulibraryapp@dpu.com>")
	email.AddTo(receiver)
	email.SetSubject("Kütüphane Uygulaması Onay Kodu")

	absPath, _ := filepath.Abs("helpers/templates/register.html")
	htmlBody, err := os.ReadFile(absPath)
	body := strings.Replace(string(htmlBody), "$code", fmt.Sprintf("%d", code), -1)

	if err != nil {
		log.Fatal(err)
	}

	email.SetBody(mail.TextHTML, body)

	err = email.Send(smtpClient)

	if err != nil {
		log.Fatal(err)
	}
}

func SendConfirmationCode(receiver string, code int) {
	server := mail.NewSMTPClient()
	server.Host = config.Config("MAIL_HOST")
	server.Port, _ = strconv.Atoi(config.Config("MAIL_PORT"))
	server.Username = config.Config("MAIL_ADDR")
	server.Password = config.Config("MAIL_PASS")

	server.Encryption = mail.EncryptionTLS

	smtpClient, err := server.Connect()
	if err != nil {
		log.Fatal(err)
	}

	email := mail.NewMSG()
	email.SetFrom("Kütüphane Uygulaması <dpulibraryapp@dpu.com>")
	email.AddTo(receiver)
	email.SetSubject("Kütüphane Uygulaması Onay Kodu")

	absPath, _ := filepath.Abs("helpers/templates/confirmation.html")
	htmlBody, err := os.ReadFile(absPath)
	body := strings.Replace(string(htmlBody), "$code", fmt.Sprintf("%d", code), -1)

	if err != nil {
		log.Fatal(err)
	}

	email.SetBody(mail.TextHTML, body)

	err = email.Send(smtpClient)

	if err != nil {
		log.Fatal(err)
	}
}
