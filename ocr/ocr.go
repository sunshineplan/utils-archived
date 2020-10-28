package ocr

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

const (
	api    = "https://api.ocr.space/parse/image"
	apikey = "5a64d478-9c89-43d8-88e3-c65de9999580"
)

type ocrText struct {
	ParsedResults []struct {
		ParsedText string
	}
}

// Read reads image from r and converts it into text
func Read(r io.Reader) (string, error) {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	defer w.Close()

	part, _ := w.CreateFormFile("file", "pic.jpg")
	if _, err := io.Copy(part, r); err != nil {
		return "", err
	}

	params := map[string]string{
		"scale": "true",
	}
	for k, v := range params {
		w.WriteField(k, v)
	}

	req, _ := http.NewRequest("POST", api, &body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("apikey", apikey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)
	var results ocrText
	if err := json.Unmarshal(b, &results); err != nil {
		return "", err
	}
	if len(results.ParsedResults) == 0 {
		return "", errors.New("no ocr result")
	}
	return results.ParsedResults[0].ParsedText, nil
}
