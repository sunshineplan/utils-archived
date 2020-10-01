package progressbar

import (
	"bytes"
	"io"
	"log"
	"os"
	"strings"
	"text/template"
	"time"
)

const defaultTemplate = ` [{{.Done}}{{.Undone}}]   {{printf "%.2f/s" .Speed}} - {{.Current -}}
{{printf "(%.2f%%)" .Percent}} of {{.Total}}   {{printf "Left: %s" .Left}}   `

// ProgressBar is a simple progress bar
type ProgressBar struct {
	width    int
	refresh  time.Duration
	template *template.Template
}

type format struct {
	Done, Undone   string
	Speed, Percent float64
	Current, Total int
	Left           time.Duration
}

func (f *format) execute(template *template.Template) {
	var buf bytes.Buffer
	if err := template.Execute(&buf, f); err == nil {
		if f.Current < f.Total {
			io.WriteString(&buf, "\r\r")
		} else {
			io.WriteString(&buf, "\r\n")
		}
	}
	io.Copy(os.Stdout, &buf)
}

// New returns a new ProgressBar with default options
func New() *ProgressBar {
	return &ProgressBar{
		width:    50,
		refresh:  5 * time.Second,
		template: template.Must(template.New("ProgressBar").Parse(defaultTemplate)),
	}
}

// SetWidth sets progress bar width
func (pb *ProgressBar) SetWidth(width int) *ProgressBar {
	pb.width = width
	return pb
}

// SetRefresh sets progress bar refresh
func (pb *ProgressBar) SetRefresh(refresh time.Duration) *ProgressBar {
	pb.refresh = refresh
	return pb
}

// SetTemplate sets progress bar template
func (pb *ProgressBar) SetTemplate(format string) *ProgressBar {
	t := template.New("ProgressBar")
	if _, err := t.Parse(format); err != nil {
		log.Print("Invalid template.")
		return pb
	}
	pb.template = t
	return pb
}

// Start ProgressBar
func (pb *ProgressBar) Start(total int, current *int) {
	var speed float64
	go func() {
		for {
			now := *current
			if now >= total {
				return
			}
			time.Sleep(pb.refresh)
			speed = float64(*current-now) / float64(pb.refresh/time.Second)
		}
	}()
	go func() {
		for {
			now := *current
			done := pb.width * now / total
			percent := float64(now) * 100 / float64(total)
			left := time.Duration(float64(total-now)/speed) * time.Second
			var progressed string
			if now < total && done != 0 {
				progressed = strings.Repeat("=", done-1) + ">"
			} else {
				progressed = strings.Repeat("=", done)
			}
			f := format{
				Done:    progressed,
				Undone:  strings.Repeat(" ", pb.width-done),
				Speed:   speed,
				Current: now,
				Percent: percent,
				Total:   total,
				Left:    left,
			}
			go f.execute(pb.template)
			if now >= total {
				return
			}
			time.Sleep(time.Second)
		}
	}()
}
