package sms

import (
	"encoding/json"
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-notifications/models"
	"github.com/WeBankPartners/wecube-plugins-notifications/services"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func SendSmsHandler(w http.ResponseWriter, r *http.Request) {
	var resp models.SmsResultObj
	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	var param models.SmsRequestObj
	err := json.Unmarshal(b, &param)
	if err != nil {
		log.Printf("Param json unmarshal error : %v \n", err)
		resp.ResultCode = "1"
		resp.ResultMessage = fmt.Sprintf("Param json unmarshal error : %v", err)
		resp.Results = models.SmsResultOutputs{Outputs: []models.SmsResultOutputObj{}}
		d, _ := json.Marshal(resp)
		w.Write(d)
		return
	}
	log.Printf("Request sms:------->> \n")
	if len(param.Inputs) == 0 {
		resp.ResultCode = "1"
		resp.ResultMessage = "Inputs list is null"
		resp.Results = models.SmsResultOutputs{Outputs: []models.SmsResultOutputObj{}}
	} else {
		resultOutputs := []models.SmsResultOutputObj{}
		for _, v := range param.Inputs {
			tmpErr := checkSmsParam(v)
			if tmpErr != nil {
				resp.ResultMessage = tmpErr.Error()
				continue
			}
			tmpErr = services.SendSms(v.SecretId, v.SecretKey, models.SmsRequest{SmsSdkAppId: v.SmsSdkAppId, PhoneNumberSet: getPhoneListFromParam(v.To), TemplateId: v.TemplateId, SignName: v.Sender, TemplateParamSet: getSmsContentParam(v.Content)})
			if tmpErr != nil {
				resp.ResultMessage = fmt.Sprintf("Send tencent cloud sms fail,%s ", tmpErr.Error())
			}
		}
		resp.Results = models.SmsResultOutputs{Outputs: resultOutputs}
		if resp.ResultMessage != "" {
			resp.ResultCode = "1"
		} else {
			resp.ResultCode = "0"
		}
	}
	d, _ := json.Marshal(resp)
	log.Printf("Result sms:------->> %s \n", string(d))
	w.Header().Set("Content-Type", "application/json")
	w.Write(d)
}

func checkSmsParam(input models.SmsInputObj) error {
	if input.SecretId == "" {
		return fmt.Errorf("Param secretId can not empty ")
	}
	if input.SecretKey == "" {
		return fmt.Errorf("Param secretKey can not empty ")
	}
	if input.TemplateId == "" {
		return fmt.Errorf("Param templateId can not empty ")
	}
	if input.SmsSdkAppId == "" {
		return fmt.Errorf("Param smsSdkAppId can not empty ")
	}
	if input.To == "" {
		return fmt.Errorf("Param to can not empty ")
	}
	if input.Content == "" {
		return fmt.Errorf("Param content can not empty ")
	}
	return nil
}

func getPhoneListFromParam(input string) []string {
	result := []string{}
	for _, v := range strings.Split(input, ",") {
		if !strings.HasPrefix(v, "+") {
			result = append(result, fmt.Sprintf("+86%s", v))
		} else {
			result = append(result, v)
		}
	}
	return result
}

func getSmsContentParam(input string) []string {
	if strings.HasSuffix(input, ";") {
		input = input[:len(input)-1]
	}
	return strings.Split(input, ";")
}
