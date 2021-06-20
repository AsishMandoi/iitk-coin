package global

import (
	"os"
)

var MaxCap = 1001.0

type Stu struct {
	Rollno   int    `json:"rollno"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Batch    string `json:"batch"`
}

type DefaultRespBodyFormat struct {
	Message interface{} `json:"message"`
	Error   interface{} `json:"error"`
}

type ViewCoinsRespBodyFormat struct {
	Message interface{} `json:"message"`
	Error   interface{} `json:"error"`
	Coins   interface{} `json:"coins"`
}

type SecretpageRespBodyFormat struct {
	Message interface{} `json:"message"`
	Error   interface{} `json:"error"`
	Data    interface{} `json:"data"`
}

type LoginRespBodyFormat struct {
	Message interface{} `json:"message"`
	Error   interface{} `json:"error"`
	Token   interface{} `json:"token"`
}

var SignatureKey = []byte(os.Getenv("SECRET_KEY"))
