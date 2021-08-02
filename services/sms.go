package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-notifications/models"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func sha256hex(s string) string {
	b := sha256.Sum256([]byte(s))
	return hex.EncodeToString(b[:])
}

func hmacsha256(s, key string) string {
	hashed := hmac.New(sha256.New, []byte(key))
	hashed.Write([]byte(s))
	return string(hashed.Sum(nil))
}

func SendSms(secretId, secretKey string, requestData models.SmsRequest) error {
	host := "sms.tencentcloudapi.com"
	algorithm := "TC3-HMAC-SHA256"
	service := "sms"
	version := "2021-01-11"
	action := "SendSms"
	region := "ap-guangzhou"
	var timestamp int64 = time.Now().Unix()

	// step 1: build canonical request string
	httpRequestMethod := "POST"
	canonicalURI := "/"
	canonicalQueryString := ""
	canonicalHeaders := "content-type:application/json; charset=utf-8\n" + "host:" + host + "\n"
	signedHeaders := "content-type;host"
	requestDataBytes, _ := json.Marshal(requestData)
	payload := string(requestDataBytes)
	hashedRequestPayload := sha256hex(payload)
	canonicalRequest := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		httpRequestMethod,
		canonicalURI,
		canonicalQueryString,
		canonicalHeaders,
		signedHeaders,
		hashedRequestPayload)
	fmt.Println(canonicalRequest)

	// step 2: build string to sign
	date := time.Unix(timestamp, 0).UTC().Format("2006-01-02")
	credentialScope := fmt.Sprintf("%s/%s/tc3_request", date, service)
	hashedCanonicalRequest := sha256hex(canonicalRequest)
	string2sign := fmt.Sprintf("%s\n%d\n%s\n%s",
		algorithm,
		timestamp,
		credentialScope,
		hashedCanonicalRequest)
	fmt.Println(string2sign)

	// step 3: sign string
	secretDate := hmacsha256(date, "TC3"+secretKey)
	secretService := hmacsha256(service, secretDate)
	secretSigning := hmacsha256("tc3_request", secretService)
	signature := hex.EncodeToString([]byte(hmacsha256(string2sign, secretSigning)))
	//fmt.Println(signature)

	// step 4: build authorization
	authorization := fmt.Sprintf("%s Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		algorithm,
		secretId,
		credentialScope,
		signedHeaders,
		signature)
	//fmt.Println(authorization)

	// step 5: send http request
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://%s", host), strings.NewReader(payload))
	if err != nil {
		return fmt.Errorf("new req error:%s ", err.Error())
	}
	req.Header.Set("Authorization", authorization)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Host", host)
	req.Header.Set("X-TC-Action", action)
	req.Header.Set("X-TC-Timestamp", fmt.Sprintf("%d", timestamp))
	req.Header.Set("X-TC-Version", version)
	req.Header.Set("X-TC-Region", region)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("http do request fail,%s ", err.Error())
	}
	fmt.Printf("response status code:%d \n", resp.StatusCode)
	var response models.SmsResponse
	respBytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return fmt.Errorf("read response body fail,%s ", err.Error())
	}
	err = json.Unmarshal(respBytes, &response)
	if err != nil {
		return fmt.Errorf("json unmarshal body fail,%s ", err.Error())
	}
	if len(response.Response.SendStatusSet) == 0 {
		return fmt.Errorf("response error:%s ", string(respBytes))
	}
	fmt.Printf("body: \n %s \n", string(respBytes))
	errorMessage := ""
	for _, v := range response.Response.SendStatusSet {
		if v.Code != "Ok" {
			errorMessage += fmt.Sprintf("phone:%s code:%s message:%s \n", v.PhoneNumber, v.Code, v.Message)
		} else {
			fmt.Printf("send sms to phone:%s success \n", v.PhoneNumber)
		}
	}
	if errorMessage != "" {
		return fmt.Errorf(errorMessage)
	}
	return nil
}
