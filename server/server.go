package server

import (
	"context"
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

// Run runs http handler.
func (o *Options) Run(handler http.Handler) error {
	server := &http.Server{Handler: handler}

	idleConnsClosed := make(chan struct{})
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		if err := server.Shutdown(context.Background()); err != nil {
			fmt.Println("Failed to close server:", err)
		}
		close(idleConnsClosed)
	}()

	if o.UNIX != "" && runtime.GOOS == "linux" {
		listener, err := net.Listen("unix", o.UNIX)
		if err != nil {
			return fmt.Errorf("Failed to listen socket file: %v", err)
		}
		// Let everyone can access the socket file
		if err := os.Chmod(o.UNIX, 0666); err != nil {
			return fmt.Errorf("Failed to chmod socket file: %v", err)
		}
		if err := server.Serve(listener); err != http.ErrServerClosed {
			return fmt.Errorf("Failed to server: %v", err)
		}
	} else {
		server.Addr = o.Host + ":" + o.Port
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			return fmt.Errorf("Failed to server: %v", err)
		}
	}
	<-idleConnsClosed
	return nil
}
