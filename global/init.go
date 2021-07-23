package global

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func init() {
	// This displays a warning if the ".env file" or ".env.dev" fails to load
	if err := godotenv.Load(".env", ".env.dev"); err != nil {
		fmt.Println("Warning: Could not load .env file; Your environment variables will be empty by default.")
	}

	SignatureKey, MyGmailId, MyPwd = []byte(os.Getenv("SECRET_KEY")), os.Getenv("GMAILID"), os.Getenv("PASSWORD")

	var parseErr error
	MaxCap, parseErr = strconv.ParseFloat(os.Getenv("MAX_CAP"), 64)
	if parseErr != nil {
		fmt.Println("Parse error for MaxCap")
	}
	MinEvents, parseErr = strconv.Atoi(os.Getenv("MIN_EVENTS"))
	if parseErr != nil {
		fmt.Println("Parse error for MinEvents")
	}
}
