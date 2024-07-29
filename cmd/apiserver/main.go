package main

import (
	"fmt"

	"github.com/MetaException/wb_l0/internal/apiserver"
	"github.com/MetaException/wb_l0/internal/cache"
	"github.com/MetaException/wb_l0/internal/natsbroker"
	"github.com/MetaException/wb_l0/internal/postgresql"
)

func main() {

	pg := postgresql.New()
	defer pg.Close()

	cs := cache.New(pg)
	cs.RestoreCache()

	ns := natsbroker.New(cs, pg)
	if err := ns.Listen(); err != nil {
		fmt.Println("nats listening error: %w", err)
	}

	apiserver := apiserver.New(cs)
	apiserver.Start()
}
