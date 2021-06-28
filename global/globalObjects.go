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

type SecretpageRespBody struct {
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

type TxnBody struct {
	Sender   int
	Receiver int
	Amount   float64
	Descr    string
}
