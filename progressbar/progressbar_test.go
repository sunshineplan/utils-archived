package progressbar

import (
	"testing"
	"time"
)

func TestProgessBar(t *testing.T) {
	var i int
	pb := New()
	go pb.Start(20, &i)
	for ; i <= 20; i++ {
		//log.Print(i)
		time.Sleep(time.Second)
	}
	time.Sleep(time.Second)
}
