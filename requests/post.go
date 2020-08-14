package requests

import "net/http"

func post(URL string, headers map[string]string, data interface{}, client *http.Client) *Response {
	req, err := buildRequest("POST", URL, data)
	if err != nil {
		return &Response{Error: err}
	}
	if headers != nil {
		addHeaders(headers, req)
	}
	return buildResponse(client.Do(req))
}

// Post request
func Post(URL string, headers map[string]string, data interface{}) *Response {
	return post(URL, headers, data, defaultClient)
}

// PostWithClient request with a custom client
func PostWithClient(URL string, headers map[string]string, data interface{}, client *http.Client) *Response {
	return post(URL, headers, data, client)
}
