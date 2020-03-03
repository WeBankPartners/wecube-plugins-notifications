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
)

type mailObj struct {
	Auth  smtp.Auth
	From  string
	Server  string
}

var (
	mailEnable  bool
	defaultMail  string
	mailAuthMap = make(map[string]mailObj)
)

func InitSmtpMail()  {
	mailEnable = true
	if !m.Config().Mail.Enable || len(m.Config().Mail.Sender) == 0 {
		mailEnable = false
		log.Println("init smtp mail fail,please check config file!")
		return
	}
	defaultMail = m.Config().Mail.Sender[0].Name
	for _,v := range m.Config().Mail.Sender {
		tmpAuth := smtp.PlainAuth(v.Token,v.User,v.Password,v.Server)
		mailAuthMap[v.Name] = mailObj{Auth:tmpAuth, From:v.User, Server:v.Server}
	}
	log.Println("init smtp mail done")
}

func sendSmtpMail(smo m.SendMailObj) error {
	if !mailEnable {
		return fmt.Errorf("mail channel is disable")
	}
	if smo.Name == "" {
		smo.Name = defaultMail
	}
	if _,b := mailAuthMap[smo.Name];!b {
		return fmt.Errorf("sender:%s is not exist", smo.Name)
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

	err := smtp.SendMail(fmt.Sprintf("%s:25", mailAuthMap[smo.Name].Server), mailAuthMap[smo.Name].Auth, mailAuthMap[smo.Name].From, smo.Accept, mailQQMessage(smo.Accept,smo.Subject,smo.Content,smo.Name,mailAuthMap[smo.Name].From))
	if err != nil {
		log.Printf("send mail error : %v \n", err)
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

func SendMailHandler(w http.ResponseWriter,r *http.Request)  {
	var resp m.MailResultObj
	b,_ := ioutil.ReadAll(r.Body)
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
	log.Printf("Send mail param: %s \n", string(b))
	if len(param.Inputs) == 0 {
		resp.ResultCode = "1"
		resp.ResultMessage = "Inputs list is null"
		resp.Results = m.MailResultOutputs{Outputs:[]m.MailResultOutputObj{}}
	}else{
		var resultOutputs []m.MailResultOutputObj
		for _,v := range param.Inputs {
			tmpResultOutputObj := m.MailResultOutputObj{CallbackParameter:v.CallbackParameter, ErrorCode:"0", ErrorMessage:""}
			cErr := sendSmtpMail(m.SendMailObj{Name:v.Sender, Accept:strings.Split(v.To, ","), Subject:v.Subject, Content:v.Content})
			if cErr != nil {
				tmpResultOutputObj.ErrorCode = "1"
				tmpResultOutputObj.ErrorMessage = fmt.Sprintf("error: %v", cErr)
				resp.ResultMessage = tmpResultOutputObj.ErrorMessage
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
	w.Write(d)
}