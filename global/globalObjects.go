package global

import "os"

type Stu struct {
	Rollno   int    `json:"rollno"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type SignupRespBodyFormat struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

type SecretpageRespBodyFormat struct {
	Message string `json:"message"`
	Error   string `json:"error"`
	Data    string `json:"data"`
}

type LoginInputFormat struct {
	Rollno   int    `json:"rollno"`
	Password string `json:"password"`
}

type LoginRespBodyFormat struct {
	Message string `json:"message"`
	Error   string `json:"error"`
	Token   string `json:"token"`
}

var SignatureKey = []byte(os.Getenv("SECRET_KEY"))
