package requests

import "net/http"

func get(URL string, headers map[string]string, client *http.Client) *Response {
	req, err := buildRequest("GET", URL, nil)
	if err != nil {
		return &Response{Error: err}
	}
	if headers != nil {
		addHeaders(headers, req)
	}
	return buildResponse(client.Do(req))
}

// Get request
func Get(URL string, headers map[string]string) *Response {
	return get(URL, headers, defaultClient)
}

// GetWithClient request with a custom client
func GetWithClient(URL string, headers map[string]string, client *http.Client) *Response {
	return get(URL, headers, client)
}
