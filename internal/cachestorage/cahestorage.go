package cachestorage

import (
	"github.com/MetaException/wb_l0/internal/postgresql"
	"github.com/MetaException/wb_l0/pkg/cache"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type CacheStorage struct {
	logger *logrus.Logger
	db     *postgresql.Postgres
	Cache  *cache.Cache
}

func New(pg *postgresql.Postgres, logger *logrus.Logger) *CacheStorage {
	return &CacheStorage{
		logger: logger,
		db:     pg,
		Cache:  cache.New(),
	}
}

func (cs *CacheStorage) RestoreCache() error {
	data, err := cs.db.GetAllData()

	if err != nil {
		cs.logger.WithError(err).Error("unable to get cache from db")
		return errors.WithStack(err)
	}

	for key, value := range data {
		cs.Cache.Set(key, value)
	}

	cs.logger.Infof("restored %v entries from db", len(data))

	return nil
}
