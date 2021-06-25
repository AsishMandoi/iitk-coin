package global

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func init() {
	// This displays a warning if .env file fails to load
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: Could not load .env file; Your environment variables will be empty by default.")
	}
	var parseErr error
	MaxCap, parseErr = strconv.ParseFloat(os.Getenv("MAX_CAP"), 64)
	if parseErr != nil {
		fmt.Println("Parse error for MaxCap")
	}
	SignatureKey = []byte(os.Getenv("SECRET_KEY"))
	// fmt.Println(MaxCap)
	// fmt.Println(SignatureKey)
}

// Defualt values
var MaxCap = 101.0
var SignatureKey = []byte("")

type Stu struct {
	Rollno   int    `json:"rollno"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Batch    string `json:"batch"`
	Role     string `json:"role"`
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
