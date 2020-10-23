package winsvc

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
)

var svcName = "Service"
var svcDescription string
var svcExecution func()

var elog debug.Log

// SetServiceName sets windows service name
func SetServiceName(name string) {
	svcName = name
}

// SetDescription sets windows service description
func SetDescription(desc string) {
	svcDescription = desc
}

// SetExecution sets windows service execution
func SetExecution(f func()) {
	svcExecution = f
}

type service struct{}

func (s *service) Execute(_ []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	changes <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}
	elog.Info(1, "Service started.")
	go svcExecution()
loop:
	for {
		c := <-r
		switch c.Cmd {
		case svc.Interrogate:
			changes <- c.CurrentStatus
			time.Sleep(100 * time.Millisecond)
			changes <- c.CurrentStatus
		case svc.Stop, svc.Shutdown:
			elog.Info(1, fmt.Sprintf("Stopping %s service(%d).", svcName, c.Context))
			break loop
		default:
			elog.Error(1, fmt.Sprintf("Unexpected control request #%d", c))
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}

func exePath() (string, error) {
	prog := os.Args[0]
	p, err := filepath.Abs(prog)
	if err != nil {
		return "", err
	}
	fi, err := os.Stat(p)
	if err == nil {
		if !fi.Mode().IsDir() {
			return p, nil
		}
		err = fmt.Errorf("%s is directory", p)
	}
	if filepath.Ext(p) == "" {
		p += ".exe"
		fi, err := os.Stat(p)
		if err == nil {
			if !fi.Mode().IsDir() {
				return p, nil
			}
			err = fmt.Errorf("%s is directory", p)
		}
	}
	return "", err
}

// InstallService installs windows service
func InstallService() error {
	exepath, err := exePath()
	if err != nil {
		return err
	}
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(svcName)
	if err == nil {
		s.Close()
		return fmt.Errorf("service %s already exists", svcName)
	}
	if svcDescription == "" {
		svcDescription = svcName
	}
	s, err = m.CreateService(svcName, exepath, mgr.Config{
		StartType:   mgr.StartAutomatic,
		Description: svcDescription,
	})
	if err != nil {
		return err
	}
	defer s.Close()
	if err := eventlog.InstallAsEventCreate(svcName, eventlog.Error|eventlog.Warning|eventlog.Info); err != nil {
		s.Delete()
		return fmt.Errorf("SetupEventLogSource() failed: %s", err)
	}
	return nil
}

// RemoveService removes windows service
func RemoveService() error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(svcName)
	if err != nil {
		return fmt.Errorf("service %s is not installed", svcName)
	}
	defer s.Close()
	if err := s.Delete(); err != nil {
		return err
	}
	if err := eventlog.Remove(svcName); err != nil {
		return fmt.Errorf("RemoveEventLogSource() failed: %s", err)
	}
	return nil
}

// RunService runs windows service
func RunService(isDebug bool) {
	var err error
	if isDebug {
		elog = debug.New(svcName)
	} else {
		elog, err = eventlog.Open(svcName)
		if err != nil {
			return
		}
	}
	defer elog.Close()

	elog.Info(1, fmt.Sprintf("Starting %s service.", svcName))
	run := svc.Run
	if isDebug {
		run = debug.Run
	}
	if err := run(svcName, &service{}); err != nil {
		elog.Error(1, fmt.Sprintf("%s service failed: %v", svcName, err))
		return
	}
	elog.Info(1, fmt.Sprintf("%s service stopped.", svcName))
}

// StartService starts windows service
func StartService() error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(svcName)
	if err != nil {
		return fmt.Errorf("could not access service: %v", err)
	}
	defer s.Close()
	if err := s.Start(); err != nil {
		return fmt.Errorf("could not start service: %v", err)
	}
	return nil
}

// StopService stops windows service
func StopService() error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(svcName)
	if err != nil {
		return fmt.Errorf("could not access service: %v", err)
	}
	defer s.Close()
	status, err := s.Control(svc.Stop)
	if err != nil {
		return fmt.Errorf("could not send control=%d: %v", svc.Stop, err)
	}
	timeout := time.Now().Add(10 * time.Second)
	for status.State != svc.Stopped {
		if timeout.Before(time.Now()) {
			return fmt.Errorf("timeout waiting for service to go to state=%d", svc.Stopped)
		}
		time.Sleep(300 * time.Millisecond)
		status, err = s.Query()
		if err != nil {
			return fmt.Errorf("could not retrieve service status: %v", err)
		}
	}
	return nil
}
