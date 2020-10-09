package progressbar

import (
	"testing"
	"time"
)

func TestProgessBar(t *testing.T) {
	done := make(chan bool, 1)
	var i int
	pb := New(done)
	go pb.Start(10, &i)
	for ; i <= 10; i++ {
		//log.Print(i)
		time.Sleep(time.Second)
	}
	<-done
}
