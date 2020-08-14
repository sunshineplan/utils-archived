package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var defaultAgent = "Chrome"
var defaultClient = &http.Client{Transport: &http.Transport{Proxy: nil}}

// SetAgent set default user agent string
func SetAgent(agent string) {
	if agent != "" {
		defaultAgent = agent
	}
}

func buildRequest(method, URL string, data interface{}) (*http.Request, error) {
	switch data.(type) {
	case nil:
		return http.NewRequest(method, URL, nil)
	case url.Values:
		req, err := http.NewRequest(method, URL, strings.NewReader(data.(url.Values).Encode()))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		return req, nil
	case interface{}, map[string]interface{}:
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		req, err := http.NewRequest(method, URL, bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		return req, nil
	default:
		return nil, fmt.Errorf("Unknown data format")
	}
}

func addHeaders(headers map[string]string, req *http.Request) {
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("User-Agent", defaultAgent)
}

// Response after request
type Response struct {
	Error      error
	Body       io.ReadCloser
	StatusCode int
	Header     http.Header
	Request    *http.Request
}

func buildResponse(resp *http.Response, err error) *Response {
	if err != nil {
		return &Response{Error: err}
	}
	return &Response{
		Error:      nil,
		Body:       resp.Body,
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		Request:    resp.Request,
	}
}

// Close response body
func (r *Response) Close() {
	if r.Error == nil {
		r.Body.Close()
	}
}

// JSON unmarshal response body to data
func (r *Response) JSON(data interface{}) error {
	if r.Error != nil {
		return r.Error
	}
	defer r.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}
	return nil
}

// Bytes return response []byte
func (r *Response) Bytes() []byte {
	if r.Error != nil {
		return nil
	}
	defer r.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil
	}
	return body
}

// String return response string
func (r *Response) String() string {
	return string(r.Bytes())
}
