package models

type SendMailObj struct {
	Name  string
	Accept  []string
	Subject  string
	Content  string
}

type MailRequestObj struct {
	RequestId  string  	`json:"requestId"`
	Inputs  []MailInputObj  `json:"inputs"`
}

type MailInputObj struct {
	CallbackParameter  string  `json:"callbackParameter"`
	Sender  string  `json:"sender"`
	To  string  `json:"to"`
	Subject  string  `json:"subject"`
	Content  string  `json:"content"`
}

type MailResultObj struct {
	ResultCode  string  `json:"resultCode"`
	ResultMessage  string  `json:"resultMessage"`
	Results  MailResultOutputs  `json:"results"`
}

type MailResultOutputs struct {
	Outputs  []MailResultOutputObj  `json:"outputs"`
}

type MailResultOutputObj struct {
	CallbackParameter  string  `json:"callbackParameter"`
	ErrorCode  string  `json:"errorCode"`
	ErrorMessage  string  `json:"errorMessage"`
}