package progressbar

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"text/template"
	"time"
)

const defaultTemplate = `[{{.Done}}{{.Undone}}]   {{printf "%.2f/s" .Speed}} - {{.Current -}}
{{printf "(%.2f%%)" .Percent}} of {{.Total}}   {{printf "Left: %s" .Left}} `

// ProgressBar is a simple progress bar.
type ProgressBar struct {
	sync.Mutex
	blockWidth int
	refresh    time.Duration
	template   *template.Template
	current    int
	total      int
	done       chan bool
	lastWidth  int
}

type format struct {
	Done, Undone   string
	Speed, Percent float64
	Current, Total int
	Left           time.Duration
}

func (f *format) execute(pb *ProgressBar) {
	var buf bytes.Buffer
	pb.template.Execute(&buf, f)
	width := buf.Len()
	if width < pb.lastWidth {
		io.WriteString(os.Stderr,
			fmt.Sprintf("\r%s\r%s", strings.Repeat(" ", pb.lastWidth), buf.Bytes()))
	} else {
		io.WriteString(os.Stderr, "\r\r"+string(buf.Bytes()))
	}
	pb.lastWidth = width
}

// New returns a new ProgressBar with default options.
func New(total int, done chan bool) *ProgressBar {
	return &ProgressBar{
		blockWidth: 50,
		refresh:    5 * time.Second,
		template:   template.Must(template.New("ProgressBar").Parse(defaultTemplate)),
		total:      total,
		done:       done,
	}
}

// SetWidth sets progress bar block width.
func (pb *ProgressBar) SetWidth(blockWidth int) *ProgressBar {
	pb.blockWidth = blockWidth
	return pb
}

// SetRefresh sets progress bar refresh time for check speed.
func (pb *ProgressBar) SetRefresh(refresh time.Duration) *ProgressBar {
	pb.refresh = refresh
	return pb
}

// SetTemplate sets progress bar template.
func (pb *ProgressBar) SetTemplate(tmplt string) error {
	t := template.New("ProgressBar")
	if _, err := t.Parse(tmplt); err != nil {
		return err
	}
	if err := t.Execute(ioutil.Discard, format{}); err != nil {
		return err
	}
	pb.template = t
	return nil
}

// Add adds the specified amount to the progress bar.
func (pb *ProgressBar) Add(num int) {
	pb.Lock()
	defer pb.Unlock()
	pb.current += num
}

// Start starts the progress bar.
func (pb *ProgressBar) Start() {
	go func() {
		start := time.Now()
		maxRefresh := pb.refresh * 3
		var speed, totalSpeed, intervalSpeed float64
		go func() {
			for {
				now := pb.current
				if now >= pb.total {
					return
				}
				time.Sleep(pb.refresh)
				totalSpeed = float64(now) / (float64(time.Since(start)) / float64(time.Second))
				intervalSpeed = float64(pb.current-now) / (float64(pb.refresh) / float64(time.Second))
				if intervalSpeed == 0 && pb.refresh < maxRefresh {
					pb.refresh += time.Second
				}
			}
		}()
		go func() {
			for {
				start := time.Now()
				now := pb.current
				if now > pb.total {
					now = pb.total
				}
				done := pb.blockWidth * now / pb.total
				percent := float64(now) * 100 / float64(pb.total)
				if intervalSpeed == 0 {
					speed = totalSpeed
				} else {
					speed = intervalSpeed
				}
				left := time.Duration(float64(pb.total-now)/speed) * time.Second
				if left < 0 {
					left = 0
				}
				var progressed string
				if now < pb.total && done != 0 {
					progressed = strings.Repeat("=", done-1) + ">"
				} else {
					progressed = strings.Repeat("=", done)
				}
				f := format{
					Done:    progressed,
					Undone:  strings.Repeat(" ", pb.blockWidth-done),
					Speed:   speed,
					Current: now,
					Percent: percent,
					Total:   pb.total,
					Left:    left,
				}
				f.execute(pb)
				if now == pb.total {
					io.WriteString(os.Stderr, "\n")
					pb.done <- true
					return
				}
				time.Sleep(time.Second - time.Since(start))
			}
		}()
	}()
}
