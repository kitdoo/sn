package signal

import (
	"os"
	"os/signal"
	"sync"
)

type Finalizer func()
type Handler func(sig os.Signal)

type Signal struct {
	mx sync.RWMutex

	signals       []os.Signal
	signalChannel chan os.Signal
	started       bool
	autoStop      bool
	stop          chan struct{}
	wg            *sync.WaitGroup
	finalizer     Finalizer
	handlers      map[os.Signal][]Handler
	bHandlers     []Handler
}

func New(signals ...os.Signal) *Signal {
	s := &Signal{
		handlers:      map[os.Signal][]Handler{},
		bHandlers:     []Handler{},
		signalChannel: make(chan os.Signal, 1),
		wg:            &sync.WaitGroup{},
		signals:       signals,
		stop:          make(chan struct{}, 1),
	}
	return s
}

func (s *Signal) SetFinalizer(f Finalizer) *Signal {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.finalizer = f
	return s
}

func (s *Signal) AddHandlerForSignal(h Handler, signals ...os.Signal) *Signal {
	s.mx.Lock()
	defer s.mx.Unlock()

	for _, sig := range signals {
		if _, ok := s.handlers[sig]; !ok {
			s.handlers[sig] = make([]Handler, 0)
		}
		s.handlers[sig] = append(s.handlers[sig], h)
	}
	return s
}

func (s *Signal) AddBroadcastHandler(h Handler) *Signal {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.bHandlers = append(s.bHandlers, h)
	return s
}

func (s *Signal) Start() *Signal {
	if s.started {
		return s
	}

	signal.Notify(s.signalChannel, s.signals...)
	s.wg.Add(1)
	go s.signalsListener()
	s.started = true
	return s
}

func (s *Signal) Stop() {
	if !s.started {
		return
	}
	s.stop <- struct{}{}
	<-s.stop
	s.started = false
}

func (s *Signal) Wait() {
	s.wg.Wait()
}

func (s *Signal) signalsListener() {
	for {
		select {
		case sig := <-s.signalChannel:
			go s.runHandlers(sig)
		case <-s.stop:
			s.runFinalizer()
			return
		}
	}
}

func (s *Signal) runHandlers(sig os.Signal) {
	s.mx.RLock()
	var handlers []Handler
	handlers = append(handlers, s.handlers[sig]...)
	handlers = append(handlers, s.bHandlers...)
	s.mx.RUnlock()

	for _, handler := range handlers {
		handler(sig)
	}
}

func (s *Signal) runFinalizer() {
	if s.finalizer != nil {
		s.finalizer()
	}
	signal.Stop(s.signalChannel)
	close(s.signalChannel)
	s.wg.Done()

	close(s.stop)
}
