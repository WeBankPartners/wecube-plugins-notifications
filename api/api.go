package api

import (
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-notifications/api/v1/mail"
	"github.com/WeBankPartners/wecube-plugins-notifications/api/v1/sms"
	"github.com/WeBankPartners/wecube-plugins-notifications/models"
	"log"
	"net/http"
)

func InitHttpServer(port int) {
	mail.InitSmtpMail()
	http.Handle("/notification/mail/send", http.HandlerFunc(mail.SendMailHandler))
	http.Handle("/notification/sms/send", http.HandlerFunc(sms.SendSmsHandler))
	listenPort := ":" + models.Config().Http.Port
	if port > 0 {
		listenPort = fmt.Sprintf(":%d", port)
	}
	log.Printf("listening %s ...\n", listenPort)
	http.ListenAndServe(listenPort, nil)
}
