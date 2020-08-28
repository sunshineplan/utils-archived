package requests

import "net/http"

// Post request
func Post(url string, headers map[string]string, data interface{}) *Response {
	return PostWithClient(url, headers, data, defaultClient)
}

// PostWithClient request with a custom client
func PostWithClient(url string, headers map[string]string, data interface{}, client *http.Client) *Response {
	return doRequest("POST", url, headers, data, client)
}
