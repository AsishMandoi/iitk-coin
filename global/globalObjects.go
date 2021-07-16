package global

// Defualt values
var SignatureKey = []byte("")
var MaxCap = 101.0
var MinEvents = 5

type Stu struct {
	Rollno   int    `json:"rollno"`
	Name     string `json:"name"`
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

type TxnBody struct {
	Sender   int
	Receiver int
	Amount   float64
	Descr    string
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
