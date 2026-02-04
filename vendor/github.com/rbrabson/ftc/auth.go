package ftc

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	username string
	authKey  string
)

// init retrieves the username and authentication key used for authentication on HTTP requests sent
// to the FTC server API endpoint.
func init() {
	godotenv.Load()

	username = os.Getenv("FTC_USERNAME")
	authKey = os.Getenv("FTC_AUTHORIZATION_KEY")
}

// SetAuthCredentials sets the username and authentication key used for authentication on HTTP requests
func SetAuthCredentials(user, key string) {
	username = user
	authKey = key
}
