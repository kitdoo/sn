// Copyright(c) 2021 Altessa Solutions Inc.
// All rights reserved.

package app

import (
	"os"

	"github.com/kitdoo/sn/internal/store/postgres"
)

func (a *App) startStoreService() (err error) {
	defer func() {
		if err != nil {
			a.logger.Error().Err(err).Msg("failed to start service")
			a.Shutdown()
			os.Exit(1)
		}
	}()

	var service *store.Store
	if service = store.New(a.config.PostgreSql, a.logger); err != nil {
		return
	}

	if err = service.Start(&a.servicesWait); err != nil {
		return
	}
	a.store = service
	a.closers = append(a.closers, a.store)
	return
}
