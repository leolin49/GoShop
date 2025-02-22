package service

import (
	"github.com/golang/glog"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"syscall"
)

type IService interface {
	Init() bool
	Reload()
	MainLoop()
	Final() bool
}

type Service struct {
	terminated bool
	Derived    IService
}

func (s *Service) Terminate() {
	s.terminated = true
}

func (s *Service) IsTerminated() bool {
	return s.terminated
}

func (s *Service) Main() bool {
	defer func() {
		if err := recover(); err != nil {
			glog.Error("[Unexcepted] ", err, "\n", string(debug.Stack()))
		}
	}()

	// catch system signal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGPIPE, syscall.SIGHUP)
	go func() {
		for sig := range ch {
			switch sig {
			case syscall.SIGHUP:
				s.Derived.Reload()
			default:
				s.Terminate()
			}
			glog.Infoln("[Service] Got signal ", sig)
		}
	}()

	runtime.GOMAXPROCS(runtime.NumCPU())

	if !s.Derived.Init() {
		return false
	}

	for !s.IsTerminated() {
		s.Derived.MainLoop()
	}

	s.Derived.Final()
	return true
}
