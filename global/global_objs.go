package global

// Defualt values
var BackendName = "__IITK-Coin.dflt-backend__"
var RedisHost = "localhost"
var RedisPassword = "__dflt-pwd__"

var MailHost = "smtp.gmail.com"
var MailPort = 587
var MyGmailId = "sender@example.com"
var MyPwd = "sender.pwd"

var SignatureKey = "xyz"
var MaxCap = 101.0
var MinEvents = 5
var TknExpTime = 15 // mins

type Stu struct {
	Rollno   int    `json:"rollno"`
	Name     string `json:"name"`
	Email    string `json:"iitk_email"`
	Password string `json:"password"`
	Batch    string `json:"batch"`
	Role     string `json:"role"`
}

type DefaultRespBody struct {
	Message interface{} `json:"message"`
	Error   interface{} `json:"error"`
}

type ViewCoinsRespBody struct {
	Message interface{} `json:"message"`
	Error   interface{} `json:"error"`
	Coins   interface{} `json:"coins"`
}

type DefaultDataRespBody struct {
	Message interface{} `json:"message"`
	Error   interface{} `json:"error"`
	Data    interface{} `json:"data"`
}

type LoginRespBody struct {
	Message interface{} `json:"message"`
	Error   interface{} `json:"error"`
	Token   interface{} `json:"token"`
}

type TxnRespBody struct {
	Message interface{} `json:"message"`
	Error   interface{} `json:"error"`
	TxnId   interface{} `json:"transaction_id"`
}

type RedeemRespBody struct {
	Message interface{} `json:"message"`
	Error   interface{} `json:"error"`
	ReqId   interface{} `json:"request_id"`
}

type TxnObj struct {
	Sender   int     `json:"sender"`
	Receiver int     `json:"receiver"`
	AmtSent  float64 `json:"amount_sent"`
	AmtRcvd  float64 `json:"amount_received"`
	Descr    string  `json:"description"`
	Otp      string  `json:"otp"`
}

type RedeemObj struct {
	Redeemer int     `json:"redeemer"`
	ItemId   int     `json:"item_id"`
	Amount   float64 `json:"amount"`
	Descr    string  `json:"description"`
	Otp      string  `json:"otp"`
}

type RedeemReqObj struct {
	Id          int         `json:"request_id"`
	Redeemer    int         `json:"redeemer"`
	ItemId      int         `json:"item_id"`
	Amount      float64     `json:"amount"`
	Descr       string      `json:"description"`
	RequestedOn interface{} `json:"requested_on"`
}

type RedeemStatusUPDBody struct {
	Id     int     `json:"request_id"`
	User   int     `json:"user"`
	Coins  float64 `json:"coins"`
	Status string  `json:"status"`
	Descr  string  `json:"description"`
}

type UserRedeemState struct {
	Id          int         `json:"id"`
	ItemId      int         `json:"item_id"`
	Amount      float64     `json:"amount"`
	Description string      `json:"description"`
	Status      string      `json:"status"`
	RequestedOn interface{} `json:"requested_on"`
	RespondedOn interface{} `json:"responded_on"`
}

type PwdResetObj struct {
	Rollno int    `json:"rollno"`
	NewPwd string `json:"new_password"`
	Otp    string `json:"otp"`
}
