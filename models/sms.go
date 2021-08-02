package models

type SmsRequest struct {
	SmsSdkAppId      string   `json:"SmsSdkAppId"`
	PhoneNumberSet   []string `json:"PhoneNumberSet"`
	TemplateId       string   `json:"TemplateId"`
	SignName         string   `json:"SignName"`
	TemplateParamSet []string `json:"TemplateParamSet"`
}

type SmsResponse struct {
	Response SmsResponseObj `json:"Response"`
}

type SmsResponseObj struct {
	RequestId     string              `json:"RequestId"`
	SendStatusSet []*SmsSendStatusObj `json:"SendStatusSet"`
}

type SmsSendStatusObj struct {
	SerialNo       string `json:"SerialNo"`
	PhoneNumber    string `json:"PhoneNumber"`
	Fee            int    `json:"Fee"`
	Code           string `json:"Code"`
	Message        string `json:"Message"`
	IsoCode        string `json:"IsoCode"`
	SessionContext string `json:"SessionContext"`
}

type SmsRequestObj struct {
	RequestId string        `json:"requestId"`
	Inputs    []SmsInputObj `json:"inputs"`
}

type SmsInputObj struct {
	CallbackParameter string `json:"callbackParameter"`
	SecretId          string `json:"secretId"`
	SecretKey         string `json:"secretKey"`
	SmsSdkAppId       string `json:"smsSdkAppId"`
	TemplateId        string `json:"templateId"`
	Sender            string `json:"sender"`
	To                string `json:"to"`
	Content           string `json:"content"`
}

type SmsResultObj struct {
	ResultCode    string           `json:"resultCode"`
	ResultMessage string           `json:"resultMessage"`
	Results       SmsResultOutputs `json:"results"`
}

type SmsResultOutputs struct {
	Outputs []SmsResultOutputObj `json:"outputs"`
}

type SmsResultOutputObj struct {
	CallbackParameter string `json:"callbackParameter"`
	ErrorCode         string `json:"errorCode"`
	ErrorMessage      string `json:"errorMessage"`
}
