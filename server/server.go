package server

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

// Options is a struct it contains unix or host and port settings.
type Options struct {
	UNIX string
	Host string
	Port string
}

// Engine is an interface of ServeHTTP and Run.
// For example, gin.Engine.
type Engine interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	Run(...string) error
}

// Run runs engine.
func (o *Options) Run(engine Engine) error {
	if o.UNIX != "" && runtime.GOOS == "linux" {
		if _, err := os.Stat(o.UNIX); err == nil {
			if err := os.Remove(o.UNIX); err != nil {
				return fmt.Errorf("Failed to remove socket file: %v", err)
			}
		}

		listener, err := net.Listen("unix", o.UNIX)
		if err != nil {
			return fmt.Errorf("Failed to listen socket file: %v", err)
		}
		if err := os.Chmod(o.UNIX, 0666); err != nil {
			return fmt.Errorf("Failed to chmod socket file: %v", err)
		}

		go func() {
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			<-quit

			if err := listener.Close(); err != nil {
				fmt.Println("Failed to close listener:", err)
			}
			if _, err := os.Stat(o.UNIX); err == nil {
				if err := os.Remove(o.UNIX); err != nil {
					fmt.Println("Failed to remove socket file:", err)
				}
			}
		}()

		return http.Serve(listener, engine)
	}
	return engine.Run(o.Host + ":" + o.Port)
}
