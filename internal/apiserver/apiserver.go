package apiserver

import (
	"io"
	"net/http"

	"github.com/MetaException/wb_l0/internal/cachestorage"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type APIServer struct {
	logger       *logrus.Logger
	router       *mux.Router
	cacheStorage *cachestorage.CacheStorage
	config       *Config
}

func New(cacheStorage *cachestorage.CacheStorage, logger *logrus.Logger, config *Config) *APIServer {
	return &APIServer{
		logger:       logger,
		router:       mux.NewRouter(),
		cacheStorage: cacheStorage,
		config:       config,
	}
}

func (s *APIServer) Start() error {
	s.logger.Info("starting api server")

	s.configureRouter()

	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *APIServer) configureRouter() {
	s.router.HandleFunc("/api/orders", s.handleGetOrderInfo())
}

func (s *APIServer) handleGetOrderInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		uid := query.Get("order_uid")

		w.Header().Set("Content-Type", "application/json")

		value, ok := s.cacheStorage.Cache.Get(uid)
		if ok {
			w.Write(value.([]byte))
		} else {
			io.WriteString(w, "order not found")
		}
	}
}
