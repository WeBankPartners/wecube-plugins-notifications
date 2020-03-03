package api

import (
	"net/http"
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-notifications/api/v1/mail"
	"github.com/WeBankPartners/wecube-plugins-notifications/models"
	"log"
)

func InitHttpServer(port int) {
	mail.InitSmtpMail()
	http.Handle("/notification/mail/send", http.HandlerFunc(mail.SendMailHandler))
	listenPort := ":" + models.Config().Http.Port
	if port > 0 {
		listenPort = fmt.Sprintf(":%d", port)
	}
	log.Printf("listening %s ...\n", listenPort)
	http.ListenAndServe(listenPort, nil)
}