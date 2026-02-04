package ftc

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	server string
)

// init retrieves the base server URL to be used when sending HTTP requests
func init() {
	godotenv.Load()

	server = os.Getenv("FTC_SERVER")
}

// SetServerURL sets the base server URL to be used when sending HTTP requests.
func SetServerURL(url string) {
	server = url
}
