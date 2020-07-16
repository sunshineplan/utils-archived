package agent

import (
	"net/http"

	"github.com/anaskhan96/soup"
)

// Get latest chrome user agent
func Get() (string, error) {
	body, err := soup.GetWithClient("https://agent.shlib.cf", &http.Client{Transport: &http.Transport{Proxy: nil}})
	if err != nil {
		return "", err
	}
	agent := soup.HTMLParse(body).Find("span", "class", "code")
	if agent.Error != nil {
		return "", agent.Error
	}
	return agent.Text(), nil
}
