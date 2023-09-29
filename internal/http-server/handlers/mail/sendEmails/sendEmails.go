package sendEmails

import (
	"api-crypto-project/internal/config"
	"api-crypto-project/internal/http-server/handlers/mail"
	"api-crypto-project/internal/http-server/handlers/mail/rate"
	"api-crypto-project/internal/lib/api/response"
	"fmt"
	"github.com/go-chi/render"
	"gopkg.in/gomail.v2"
	"log"
	"net/http"
	"strings"
)

const op = "handlers.mail.sendMails.New"

type MAILSender interface {
	SendEmails(cryptoCredits mail.Credits) error
}

func New(mailSender MAILSender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		crypto, err := getCredits(w, r)
		if err != nil {
			return
		}

		if err = mailSender.SendEmails(crypto); err != nil {
			log.Print(err)
			render.JSON(w, r, response.Error(op+":failed to send mails"))
			return
		}

		render.JSON(w, r, response.OK())
	}
}

func SendToSingleMail(mailToSend string, cryptoCredits mail.Credits) error {
	cfg := config.MustLoad()
	msg := gomail.NewMessage()

	msg.SetHeader("From", cfg.MailSender)
	msg.SetHeader("To", mailToSend)
	msg.SetHeader("Subject", "Current "+cryptoCredits.ID+" currency!")
	msg.SetBody("text/plain", "Current currency of "+cryptoCredits.ID+" is "+cryptoCredits.PRICE+" "+cryptoCredits.CURRENCY)

	d := gomail.NewDialer("smtp.gmail.com", 587, cfg.MailSender, cfg.AppPassword)
	if err := d.DialAndSend(msg); err != nil {
		return fmt.Errorf("error sending")
	}

	return nil
}

func getCredits(w http.ResponseWriter, r *http.Request) (mail.Credits, error) {
	var crypto mail.Crypto
	const op = "handlers.mail.sendMails.getCredits"

	resp, err := rate.RequestProcessing(w, r)
	if err != nil {
		return crypto.Credits, err
	}

	respBody, err := rate.ResponseRead(w, r, resp)
	if err != nil {
		return crypto.Credits, fmt.Errorf(op + ": failed to read response")
	}

	crypto, err = rate.UnmarshalCredits(respBody)
	if err != nil {
		return crypto.Credits, fmt.Errorf(op + ": failed to unmarshal credits")
	}

	price := strings.Split(crypto.Credits.PRICE, ".")
	crypto.Credits.PRICE = price[0]

	return crypto.Credits, nil
}
