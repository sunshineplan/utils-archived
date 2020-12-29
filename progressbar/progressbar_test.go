package progressbar

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestProgessBar(t *testing.T) {
	pb := New(15).SetRefresh(4 * time.Second)
	pb.Start()
	for i := 0; i < pb.total; i++ {
		//log.Print(i)
		pb.Add(1)
		time.Sleep(time.Second)
	}
	<-pb.Done
	pb = New(10).SetRefresh(500 * time.Millisecond)
	pb.Start()
	for i := 0; i < pb.total; i++ {
		//log.Print(i)
		pb.Add(1)
		time.Sleep(time.Second)
	}
	<-pb.Done
}

func TestFromReader(t *testing.T) {
	resp, err := http.Get("https://golang.org/dl/go1.15.6.windows-amd64.msi")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	total, err := strconv.Atoi(resp.Header.Get("content-length"))
	if err != nil {
		t.Fatal(err)
	}
	pb := New(total).SetUnit("bytes")
	if err := pb.FromReader(resp.Body, ioutil.Discard); err != nil {
		t.Fatal(err)
	}
	<-pb.Done
}

func TestSetTemplate(t *testing.T) {
	pb := &ProgressBar{}
	if err := pb.SetTemplate(`{{.Done}}`); err != nil {
		t.Error(err)
	}
	if err := pb.SetTemplate(`{{.Test}}`); err == nil {
		t.Error("gave nil error; want error")
	}
}
