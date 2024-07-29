package cache

import (
	"github.com/MetaException/wb_l0/internal/postgresql"
	"github.com/sirupsen/logrus"
)

type CacheStorage struct {
	logger *logrus.Logger
	db     *postgresql.Postgres
}

var storage map[string]any

func New(pg *postgresql.Postgres) *CacheStorage {
	storage = make(map[string]any)
	return &CacheStorage{
		logger: logrus.New(),
		db:     pg,
	}
}

func (cs *CacheStorage) Set(uid string, data any) {
	storage[uid] = data
}

func (cs *CacheStorage) Get(uid string) (any, bool) {

	value, ok := storage[uid]

	if !ok {
		return nil, false
	}

	return value, true
}

func (cs *CacheStorage) RestoreCache() {
	data, _ := cs.db.GetAllData()

	for key, value := range data {
		cs.Set(key, value)
	}
}
