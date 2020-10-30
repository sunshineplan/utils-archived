package progressbar

import (
	"testing"
	"time"
)

func TestProgessBar(t *testing.T) {
	done := make(chan bool, 1)
	pb := New(15, done).SetRefresh(4 * time.Second)
	pb.Start()
	for i := 0; i < pb.total; i++ {
		//log.Print(i)
		pb.Add(1)
		time.Sleep(time.Second)
	}
	<-done
	pb = New(10, done).SetRefresh(500 * time.Millisecond)
	pb.Start()
	for i := 0; i < pb.total; i++ {
		//log.Print(i)
		pb.Add(1)
		time.Sleep(time.Second)
	}
	<-done
}

func TestSetTemplate(t *testing.T) {
	pb := ProgressBar{}
	if err := pb.SetTemplate(`{{.Done}}`); err != nil {
		t.Error("test except non error")
	}
	if err := pb.SetTemplate(`{{.Test}}`); err == nil {
		t.Error("test except error")
	}
}
