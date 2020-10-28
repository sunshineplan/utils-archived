package requests

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetAndHead(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, world!")
	}))
	defer ts.Close()

	resp := Get(ts.URL, H{"hello": "world"})
	if resp.Error != nil {
		t.Error(resp.Error)
	}
	if resp.Request.Method != "GET" {
		t.Errorf("expected method %q; got %q", "GET", resp.Request.Method)
	}
	if resp.Request.URL.String() != ts.URL {
		t.Errorf("expected URL %q; got %q", ts.URL, resp.Request.URL.String())
	}
	if h := resp.Request.Header.Get("hello"); h != "world" {
		t.Errorf("expected hello header %q; got %q", "world", h)
	}
	if ua := resp.Request.Header.Get("user-agent"); ua != "Chrome" {
		t.Errorf("expected user agent %q; got %q", "Chrome", ua)
	}
	if s := resp.String(); s != "Hello, world!" {
		t.Error("Incorrect get response body:", s)
	}

	resp = Head(ts.URL, H{"hello": "world"})
	if resp.Error != nil {
		t.Error(resp.Error)
	}
	if resp.Request.Method != "HEAD" {
		t.Errorf("expected method %q; got %q", "HEAD", resp.Request.Method)
	}
	if l := resp.Request.ContentLength; l != 0 {
		t.Error("Incorrect head response body:", l)
	}
}

func TestPost(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := ioutil.ReadAll(r.Body)
		fmt.Fprint(w, string(c))
	}))
	defer ts.Close()

	SetAgent("test")
	resp := Post(ts.URL, H{"hello": "world"}, nil)
	if resp.Error != nil {
		t.Error(resp.Error)
	}
	defer resp.Close()
	if resp.Request.Method != "POST" {
		t.Errorf("expected method %q; got %q", "POST", resp.Request.Method)
	}
	if resp.Request.URL.String() != ts.URL {
		t.Errorf("expected URL %q; got %q", ts.URL, resp.Request.URL.String())
	}
	if h := resp.Request.Header.Get("hello"); h != "world" {
		t.Errorf("expected hello header %q; got %q", "world", h)
	}
	if ua := resp.Request.Header.Get("user-agent"); ua != "test" {
		t.Errorf("expected user agent %q; got %q", "test", ua)
	}
	if l := resp.Request.ContentLength; l != 0 {
		t.Error("Incorrect response body:", l)
	}

	resp = Post(ts.URL, nil, url.Values{"test": []string{"test"}})
	if resp.Error != nil {
		t.Error(resp.Error)
	}
	defer resp.Close()
	if ct := resp.Request.Header.Get("Content-Type"); ct != "application/x-www-form-urlencoded" {
		t.Errorf("expected Content-Type header %q; got %q", "application/x-www-form-urlencoded", ct)
	}
	if s := resp.String(); s != "test=test" {
		t.Error("Incorrect response body:", s)
	}

	resp = Post(ts.URL, nil, map[string]interface{}{"test": "test"})
	if resp.Error != nil {
		t.Error(resp.Error)
	}
	defer resp.Close()
	if ct := resp.Request.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type header %q; got %q", "application/json", ct)
	}
	var json struct{ Test string }
	if err := resp.JSON(&json); err != nil {
		t.Error(err)
	}
	if json != struct{ Test string }{Test: "test"} {
		t.Error("Incorrect response body:", json)
	}
}
