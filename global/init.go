package global

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func init() {
	// This displays a warning if the ".env" file fails to load.
	if err := godotenv.Load(); err != nil {
		fmt.Println("WARNING: .env file could not be loaded")
	}

	BackendName, RedisHost, RedisPassword = os.Getenv("BACKEND_CONTAINER_NAME"), os.Getenv("REDIS_CONTAINER_NAME"), os.Getenv("REDIS_PWD")
	SignatureKey, MailHost, MyGmailId, MyPwd = os.Getenv("SECRET_KEY"), os.Getenv("MAIL_HOST"), os.Getenv("EMAIL_ID"), os.Getenv("PASSWORD")

	var parseErr error
	MailPort, parseErr = strconv.Atoi(os.Getenv("MAIL_PORT"))
	if parseErr != nil {
		fmt.Println("Parse error for MailPort")
	}
	MaxCap, parseErr = strconv.ParseFloat(os.Getenv("MAX_CAP"), 64)
	if parseErr != nil {
		fmt.Println("Parse error for MaxCap")
	}
	MinEvents, parseErr = strconv.Atoi(os.Getenv("MIN_EVENTS"))
	if parseErr != nil {
		fmt.Println("Parse error for MinEvents")
	}
	TknExpTime, parseErr = strconv.Atoi(os.Getenv("EXP_TIME"))
	if parseErr != nil {
		fmt.Println("Parse error for TknExpTime")
	}
}
