package main

import (
	"github.com/MetaException/wb_l0/internal/apiserver"
	"github.com/MetaException/wb_l0/internal/cache"
	"github.com/MetaException/wb_l0/internal/natsbroker"
	"github.com/MetaException/wb_l0/internal/postgresql"
	"github.com/sirupsen/logrus"
)

func main() {

	pg := postgresql.New()
	defer pg.Close()

	cs := cache.New(pg)
	if err := cs.RestoreCache(); err != nil {
		logrus.WithError(err).Error("unable to restore cache from db")
	}

	ns := natsbroker.New(cs, pg)
	defer ns.Close()

	if err := ns.Listen(); err != nil {
		logrus.WithError(err).Error("nats listening error")
	}

	apiserver := apiserver.New(cs)
	apiserver.Start()
}
