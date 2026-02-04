package ftc

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
)

// Get sends an HTTP GET request to the FTC Server API endpoint and returns the data provided by
// that endpoint in a byte array.
func getURL(url string) ([]byte, error) {
	// Setup the HTTP client for the request
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// Set up HTTP request with basic authorization.
	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(username, authKey)

	// Send the request and get the response
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// HTTP request was not successful
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		err = fmt.Errorf("HTTP Status Code: %d (%s)", resp.StatusCode, http.StatusText(resp.StatusCode))
		return nil, err
	}

	// Read the output from the server
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
