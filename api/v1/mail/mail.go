package mail

import (
	"net/smtp"
	"fmt"
	"bytes"
	"strings"
	m "github.com/WeBankPartners/wecube-plugins-notifications/models"
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"crypto/tls"
	"regexp"
)

type mailObj struct {
	Auth  smtp.Auth
	IsSSL  bool
	From  string
	Server  string
}

var (
	mailEnable  bool
	defaultMail  string
	//mailAuthMap = make(map[string]mailObj)
	defaultAuth  smtp.Auth
	defaultServer  string
	defaultFrom  string
)

func InitSmtpMail()  {
	mailEnable = true
	if !m.Config().Mail.Enable || len(m.Config().Mail.Sender) == 0 {
		mailEnable = false
		log.Println("init smtp mail fail,please check config file!")
		return
	}
	defaultMail = m.Config().Mail.Sender[0].Name
	for i,v := range m.Config().Mail.Sender {
		if v.Server == "" || v.Server == "default_server" || i > 0 {
			continue
		}
		if v.Password != "" {
			tmpPassword := m.DecryptRsa(v.Password)
			defaultAuth = smtp.PlainAuth(v.Token, v.User, tmpPassword, v.Server)
		}
		//mailAuthMap[v.Name] = mailObj{Auth:tmpAuth, From:v.User, Server:v.Server}
		defaultServer = v.Server
		defaultFrom = v.User
	}
	log.Println("init smtp mail done")
}

func sendSMTPMail(smo m.SendMailObj) error {
	if !mailEnable {
		return fmt.Errorf("mail channel is disable")
	}
	if smo.Name == "" {
		smo.Name = defaultMail
	}
	if len(smo.Accept) == 0 {
		return fmt.Errorf("param to is null")
	}
	if smo.Subject == "" {
		return fmt.Errorf("subject is empty")
	}
	if smo.Content == "" {
		return fmt.Errorf("content is empty")
	}
	var err error
	tmpAuth := defaultAuth
	tmpServer := defaultServer
	tmpFrom := defaultFrom
	if smo.Sender != "" && smo.SenderServer != "" {
		if smo.SenderPassword != "" {
			tmpAuth = smtp.PlainAuth("", smo.Sender, smo.SenderPassword, smo.SenderServer)
		}else{
			tmpAuth = nil
		}
		tmpServer = smo.SenderServer
		tmpFrom = smo.Sender
		log.Printf("use param server:%s user:%s pw:%s \n", smo.SenderServer, smo.Sender, smo.SenderPassword)
	}else{
		if defaultServer == "" || defaultFrom == "" {
			return fmt.Errorf("param sender server is empty and default config server is empty,no specify server")
		}
		log.Println("use default config mail")
	}
	if smo.SSL {
		err = sendSMTPMailTLS(smo, tmpAuth, tmpServer, tmpFrom)
	}else {
		var address string
		if strings.Contains(tmpServer, ":") {
			address = tmpServer
		}else {
			address = fmt.Sprintf("%s:25", tmpServer)
		}
		err = smtp.SendMail(address, tmpAuth, tmpFrom, smo.Accept, mailQQMessage(smo.Accept, smo.Subject, smo.Content, smo.Name, tmpFrom))
	}
	return err
}

func sendSMTPMailTLS(smo m.SendMailObj,tmpAuth smtp.Auth,tmpServer,tmpFrom string) error {
	tlsConfig := &tls.Config{
		InsecureSkipVerify:true,
		ServerName: tmpServer,
	}
	var address string
	if strings.Contains(tmpServer, ":") {
		address = tmpServer
	}else {
		address = fmt.Sprintf("%s:465", tmpServer)
	}
	conn,err := tls.Dial("tcp", address, tlsConfig)
	if err != nil {
		return fmt.Errorf("tls dial error: %v", err)
	}
	client,err := smtp.NewClient(conn, tmpServer)
	if err != nil {
		return fmt.Errorf("smtp new client error: %v", err)
	}
	defer client.Close()
	if b,_ := client.Extension("AUTH"); b {
		err = client.Auth(tmpAuth)
		if err != nil {
			return fmt.Errorf("client auth error: %v", err)
		}
	}
	err = client.Mail(tmpFrom)
	if err != nil {
		return fmt.Errorf("client mail set from error: %v", err)
	}
	for _,to := range smo.Accept {
		if err = client.Rcpt(to); err != nil {
			return fmt.Errorf("client rcpt %s error: %v", to, err)
		}
	}
	w,err := client.Data()
	if err != nil {
		return fmt.Errorf("client data init error: %v", err)
	}
	_,err = w.Write(mailQQMessage(smo.Accept, smo.Subject, smo.Content, smo.Name, tmpFrom))
	if err != nil {
		return fmt.Errorf("write message error: %v", err)
	}
	w.Close()
	err = client.Quit()
	if err != nil {
		return fmt.Errorf("client quit error: %v", err)
	}
	return err
}

func mailQQMessage(to []string,subject,content,sender,sendFrom string) []byte {
	var buff bytes.Buffer
	buff.WriteString("To:")
	buff.WriteString(strings.Join(to, ","))
	buff.WriteString("\r\nFrom:")
	buff.WriteString(sender+"<"+sendFrom+">")
	buff.WriteString("\r\nSubject:")
	buff.WriteString(subject)
	buff.WriteString("\r\nContent-Type:text/plain;charset=UTF-8\r\n\r\n")
	buff.WriteString(content)
	return buff.Bytes()
}

func verifyMailAddress(mailString string) bool {
	reg := regexp.MustCompile(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`)
	return reg.MatchString(mailString)
}

func SendMailHandler(w http.ResponseWriter,r *http.Request)  {
	var resp m.MailResultObj
	b,_ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	var param m.MailRequestObj
	err := json.Unmarshal(b, &param)
	if err != nil {
		log.Printf("Param json unmarshal error : %v \n", err)
		resp.ResultCode = "1"
		resp.ResultMessage = fmt.Sprintf("Param json unmarshal error : %v",err)
		resp.Results = m.MailResultOutputs{Outputs:[]m.MailResultOutputObj{}}
		d,_ := json.Marshal(resp)
		w.Write(d)
		return
	}
	log.Printf("Request:------->> %s \n", string(b))
	if len(param.Inputs) == 0 {
		resp.ResultCode = "1"
		resp.ResultMessage = "Inputs list is null"
		resp.Results = m.MailResultOutputs{Outputs:[]m.MailResultOutputObj{}}
	}else{
		var resultOutputs []m.MailResultOutputObj
		for _,v := range param.Inputs {
			tmpResultOutputObj := m.MailResultOutputObj{CallbackParameter:v.CallbackParameter, ErrorCode:"0", ErrorMessage:""}
			isSSl := false
			if v.SSL == "Y" || v.SSL == "y" {
				isSSl = true
			}
			v.To = strings.Replace(v.To, "[", "", -1)
			v.To = strings.Replace(v.To, "]", "", -1)
			toList := strings.Split(v.To, ",")
			for _,vv := range toList {
				if !verifyMailAddress(vv) {
					log.Printf("Index: %s ,mail: %s validate fail", v.CallbackParameter, vv)
					tmpResultOutputObj.ErrorCode = "1"
					tmpResultOutputObj.ErrorMessage = fmt.Sprintf("Index: %s ,mail: %s validate fail", v.CallbackParameter, vv)
					resp.ResultMessage = tmpResultOutputObj.ErrorMessage
					break
				}
			}
			if tmpResultOutputObj.ErrorCode == "0" {
				v.SenderPassword = m.DecryptRsa(v.SenderPassword)
				cErr := sendSMTPMail(m.SendMailObj{Name: v.SenderMail, Accept: toList, Subject: v.Subject, Content: v.Content, SSL: isSSl, Sender:v.SenderMail, SenderPassword:v.SenderPassword, SenderServer:v.SenderMailServer})
				if cErr != nil {
					log.Printf("Index: %s ,send mail error : %v \n", v.CallbackParameter, cErr)
					tmpResultOutputObj.ErrorCode = "1"
					tmpResultOutputObj.ErrorMessage = fmt.Sprintf("error: %v", cErr)
					resp.ResultMessage = tmpResultOutputObj.ErrorMessage
				}
			}
			resultOutputs = append(resultOutputs, tmpResultOutputObj)
		}
		resp.Results = m.MailResultOutputs{Outputs:resultOutputs}
		if resp.ResultMessage != "" {
			resp.ResultCode = "1"
		}else{
			resp.ResultCode = "0"
		}
	}
	d,_ := json.Marshal(resp)
	log.Printf("Result:------->> %s \n", string(d))
	w.Header().Set("Content-Type", "application/json")
	w.Write(d)
}