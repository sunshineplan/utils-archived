package progressbar

import (
	"testing"
	"time"
)

func TestProgessBar(t *testing.T) {
	done := make(chan bool, 1)
	var i int
	pb := New(done)
	go pb.Start(20, &i)
	for ; i <= 20; i++ {
		//log.Print(i)
		time.Sleep(time.Second)
	}
	<-done
}
