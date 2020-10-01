package progressbar

import (
	"testing"
	"time"
)

func TestProgessBar(t *testing.T) {
	var i int
	pb := New()
	go pb.Start(20, &i)
	for i = 0; i < 20; i++ {
		//log.Print(i + 1)
		time.Sleep(time.Second)
	}
	time.Sleep(time.Second * 3)
}
