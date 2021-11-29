package app

import (
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/kitdoo/sn/internal/signal"

	"github.com/kitdoo/sn/internal/atomic"
	"github.com/kitdoo/sn/internal/path"
	"github.com/kitdoo/sn/internal/store/postgres"

	"github.com/kitdoo/sn/internal/config"
	grpcService "github.com/kitdoo/sn/internal/transport/grpc"
	"github.com/rs/zerolog"
)

type InternalService interface {
	Shutdown() error
}

type App struct {
	servicesWait sync.WaitGroup
	configPath   string
	config       *config.Config
	signal       *signal.Signal
	logger       zerolog.Logger
	grpc         *grpcService.Server
	store        *store.Store
	started      atomic.Bool
	closers      []InternalService
}

func New(configPath string) *App {
	s := &App{
		configPath: configPath,
	}
	return s
}

func (a *App) Start() error {
	if !a.started.CompareAndSwap(false, true) {
		return nil
	}

	cfg, err := config.Load(a.configPath)
	if err != nil {
		return err
	}
	a.config = cfg

	a.loggerSetup()

	_ = os.MkdirAll(path.VarDir(), os.ModePerm)
	_ = os.MkdirAll(path.LibDir(), os.ModePerm)

	a.logger.Info().Msgf("bin dir: %s", path.BinDir())
	a.logger.Info().Msgf("etc dir: %s", path.EtcDir())
	a.logger.Info().Msgf("var dir: %s", path.VarDir())
	a.logger.Info().Msgf("lib dir: %s", path.LibDir())

	a.signal = signal.New(
		syscall.SIGINT,  // TERMINAL INTERRUPT
		syscall.SIGTERM, // TERMINATION
		syscall.SIGQUIT, // QUIT
		syscall.SIGABRT, // ABORT
	).Start()

	a.signal.AddBroadcastHandler(func(sig os.Signal) {
		a.logger.Info().Str("signal", sig.String()).Msg("received signal")
		if sig == syscall.SIGINT || sig == syscall.SIGTERM || sig == syscall.SIGQUIT {
			a.signal.Stop()
		}
	})

	exitOnErr := func(f func() error) {
		if err == nil {
			err = f()
		}
	}

	exitOnErr(a.startStoreService)
	exitOnErr(a.startGRPCService)

	if err != nil {
		return err
	}
	return nil
}

func (a *App) Wait() *App {
	a.signal.Wait()
	return a
}

func (a *App) shutdown() {
	a.signal.Stop()
	for _, s := range a.closers {
		_ = s.Shutdown()
	}
	a.servicesWait.Wait()
}

const (
	tickerTime  = time.Second * 15
	waitingTime = time.Minute * 2
)

func (a *App) Shutdown() {
	ch := make(chan struct{}, 1)

	ticker := time.NewTicker(tickerTime)
	defer ticker.Stop()

	timer := time.NewTimer(waitingTime) //nolint:gomnd
	defer timer.Stop()

	go func() {
		a.shutdown()
		ch <- struct{}{}
	}()

	select {
	case <-ticker.C:
		a.logger.Warn().
			Dur("after", tickerTime). //nolint:gomnd
			Msg("process still running")
	case <-timer.C:
		a.logger.Warn().Msg("process killed")
		syscall.Exit(0)
	case <-ch:
		break
	}

	a.logger.Info().Msg("process have been stopped. Good Bay.")
	time.Sleep(time.Second)
}
