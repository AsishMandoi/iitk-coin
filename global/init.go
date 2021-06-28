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
	SignatureKey = []byte(os.Getenv("SECRET_KEY"))
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
