package store

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/kitdoo/sn/internal/atomic"
	"github.com/kitdoo/sn/internal/config"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

var (
	ErrConnection = errors.New("failed to establish connection")
)

type Store struct {
	config  *config.PostgreSql
	logger  zerolog.Logger
	db      *sql.DB
	wait    *sync.WaitGroup
	started atomic.Bool
}

func New(config *config.PostgreSql, logger zerolog.Logger) *Store {
	m := &Store{
		config: config,
		logger: logger.With().Str("subsystem", "store:postgres").Logger(),
	}
	return m
}

func (m *Store) Start(wait *sync.WaitGroup) error {
	if !m.started.CompareAndSwap(false, true) {
		return nil
	}

	m.wait = wait
	m.wait.Add(1)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		m.config.Host, m.config.Port, m.config.Username, m.config.Password, m.config.Database)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		m.logger.Error().Err(err).
			Str("database", m.config.Database).
			Msg("failed to establish connection")
		return nil
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("%w: %s", ErrConnection, err.Error())
	}
	m.db = db
	m.logger.Info().
		Str("database", m.config.Database).
		Msg("connection was successfully established")

	return nil
}

func (m *Store) Shutdown() error {
	if !m.started.CompareAndSwap(true, false) {
		return nil
	}

	defer func() {
		m.logger.Info().Msg("service was stopped")
		m.wait.Done()
	}()

	if err := m.db.Close(); err != nil {
		return fmt.Errorf("%s", err.Error())
	}

	return nil
}
