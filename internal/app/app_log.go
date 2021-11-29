package app

import (
	"fmt"
	stdlog "log"
	"os"
	"time"

	"github.com/kitdoo/sn/internal/version"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
)

func (a *App) loggerSetup() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	level, err := zerolog.ParseLevel(a.config.Log.Level)
	if err != nil {
		level = zerolog.ErrorLevel
	}
	zerolog.SetGlobalLevel(level)
	zerolog.DisableSampling(true)

	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		NoColor:    false,
	}

	//nolint:gomnd
	wr := diode.NewWriter(output, 1000, 30*time.Millisecond, func(missed int) {
		fmt.Printf("logger dropped %d messages", missed)
	})

	l := zerolog.New(wr).With().Timestamp().
		Str("app_name", version.AppName).
		Str("app_version", version.Version)
	for k, v := range a.config.Log.Tags {
		l = l.Str(k, v)
	}
	a.logger = l.Logger()

	sl := a.logger.With().Logger().Level(zerolog.InfoLevel)
	stdlog.SetOutput(sl)
	stdlog.SetFlags(0)
}
