package requests

import "net/http"

// Get request
func Get(url string, headers map[string]string) *Response {
	return GetWithClient(url, headers, defaultClient)
}

// GetWithClient request with a custom client
func GetWithClient(url string, headers map[string]string, client *http.Client) *Response {
	return doRequest("GET", url, headers, nil, client)
}
