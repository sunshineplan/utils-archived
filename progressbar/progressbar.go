package progressbar

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"
)

const defaultTemplate = `[{{.Done}}{{.Undone}}]   {{.Speed}} - {{.Current -}}
({{.Percent}}) of {{.Total}}   {{printf "Left: %s" .Left}} `

// ProgressBar is a simple progress bar.
type ProgressBar struct {
	sync.Mutex
	blockWidth int
	refresh    time.Duration
	template   *template.Template
	current    int
	total      int
	Done       chan bool
	lastWidth  int
	speed      float64
	unit       string
}

type counter struct{ *ProgressBar }

func (c *counter) Write(b []byte) (int, error) {
	c.Add(len(b))
	return 0, nil
}

type format struct {
	Done, Undone   string
	Speed, Percent string
	Current, Total string
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

func humanizeBytes(n int) string {
	if n < 10 {
		return fmt.Sprintf("%dB", n)
	}
	e := math.Floor(math.Log(float64(n)) / math.Log(1000))
	suffix := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}[int(e)]
	val := math.Floor(float64(n)/math.Pow(1000, e)*10+0.5) / 10
	format := "%.0f%s"
	if val < 100 {
		format = "%.1f%s"
	}
	return fmt.Sprintf(format, val, suffix)
}

// New returns a new ProgressBar with default options.
func New(total int) *ProgressBar {
	return &ProgressBar{
		blockWidth: 50,
		refresh:    5 * time.Second,
		template:   template.Must(template.New("ProgressBar").Parse(defaultTemplate)),
		total:      total,
		Done:       make(chan bool, 1),
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

// SetUnit sets progress bar unit.
func (pb *ProgressBar) SetUnit(unit string) *ProgressBar {
	pb.unit = unit
	return pb
}

// Add adds the specified amount to the progress bar.
func (pb *ProgressBar) Add(num int) {
	pb.Lock()
	defer pb.Unlock()
	pb.current += num
}

func (pb *ProgressBar) startRefresh() {
	start := time.Now()
	maxRefresh := pb.refresh * 3
	for {
		now := pb.current
		if now >= pb.total {
			return
		}
		time.Sleep(pb.refresh)
		totalSpeed := float64(now) / (float64(time.Since(start)) / float64(time.Second))
		intervalSpeed := float64(pb.current-now) / (float64(pb.refresh) / float64(time.Second))
		if intervalSpeed == 0 {
			pb.speed = totalSpeed
		} else {
			pb.speed = intervalSpeed
		}
		if intervalSpeed == 0 && pb.refresh < maxRefresh {
			pb.refresh += time.Second
		}
	}
}

func (pb *ProgressBar) startCount() {
	for {
		start := time.Now()
		now := pb.current
		if now > pb.total {
			now = pb.total
		}
		done := pb.blockWidth * now / pb.total
		percent := float64(now) * 100 / float64(pb.total)

		left := time.Duration(float64(pb.total-now)/pb.speed) * time.Second
		if left < 0 {
			left = 0
		}
		var progressed string
		if now < pb.total && done != 0 {
			progressed = strings.Repeat("=", done-1) + ">"
		} else {
			progressed = strings.Repeat("=", done)
		}
		var f format
		if pb.unit == "bytes" {
			f = format{
				Done:    progressed,
				Undone:  strings.Repeat(" ", pb.blockWidth-done),
				Speed:   humanizeBytes(int(pb.speed)) + "/s",
				Current: humanizeBytes(now),
				Percent: fmt.Sprintf("%.2f%%", percent),
				Total:   humanizeBytes(pb.total),
				Left:    left,
			}
		} else {
			f = format{
				Done:    progressed,
				Undone:  strings.Repeat(" ", pb.blockWidth-done),
				Speed:   fmt.Sprintf("%.2f/s", pb.speed),
				Current: strconv.Itoa(now),
				Percent: fmt.Sprintf("%.2f%%", percent),
				Total:   strconv.Itoa(pb.total),
				Left:    left,
			}
		}
		f.execute(pb)
		if now == pb.total {
			io.WriteString(os.Stderr, "\n")
			pb.Done <- true
			return
		}
		time.Sleep(time.Second - time.Since(start))
	}
}

// Start starts the progress bar.
func (pb *ProgressBar) Start() {
	go func() {
		go pb.startRefresh()
		go pb.startCount()
	}()
}

// FromReader starts the progress bar from a reader.
func (pb *ProgressBar) FromReader(r io.Reader, w io.Writer) (written int64, err error) {
	go func() {
		go pb.startRefresh()
		go pb.startCount()
	}()
	if written, err = io.Copy(w, io.TeeReader(r, &counter{pb})); err != nil {
		return
	}
	return
}
