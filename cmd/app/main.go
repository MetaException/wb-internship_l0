package main

import (
	"github.com/MetaException/wb_l0/internal/apiserver"
	"github.com/MetaException/wb_l0/internal/cachestorage"
	"github.com/MetaException/wb_l0/internal/natsbroker"
	"github.com/MetaException/wb_l0/internal/postgresql"
	"github.com/MetaException/wb_l0/pkg/logger"
)

func main() {

	logger := logger.NewLogrus()

	pgConfig := postgresql.NewEnvConfig()
	pg := postgresql.New(logger, pgConfig)
	defer pg.Close()

	cs := cachestorage.New(pg, logger)
	if err := cs.RestoreCache(); err != nil {
		logger.WithError(err).Error("unable to restore cache from db")
	}

	natsConfig := natsbroker.NewEnvConfig()
	ns := natsbroker.New(cs, pg, logger, natsConfig)
	defer ns.Close()

	if err := ns.Listen(); err != nil {
		logger.WithError(err).Error("nats listening error")
	}

	apiServerConfig := apiserver.NewEnvConfig()
	apiserver := apiserver.New(cs, logger, apiServerConfig)
	apiserver.Start()
}
