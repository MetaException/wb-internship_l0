package apiserver

import (
	"io"
	"net/http"

	"github.com/MetaException/wb_l0/internal/cache"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type APIServer struct {
	logger       *logrus.Logger
	router       *mux.Router
	cacheStorage *cache.CacheStorage
}

func New(cacheStorage *cache.CacheStorage) *APIServer {
	return &APIServer{
		logger:       logrus.New(),
		router:       mux.NewRouter(),
		cacheStorage: cacheStorage,
	}
}

func (s *APIServer) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	s.logger.Info("starting api server")

	s.configureRouter()

	return http.ListenAndServe(":8080", s.router)
}

func (s *APIServer) configureLogger() error {
	level, err := logrus.ParseLevel("debug")
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)

	return nil
}

func (s *APIServer) configureRouter() {
	s.router.HandleFunc("/api/orders", s.handleGetOrderInfo())
}

func (s *APIServer) handleGetOrderInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		uid := query.Get("order_uid")

		w.Header().Set("Content-Type", "application/json")

		value, ok := s.cacheStorage.Get(uid)
		if ok {
			w.Write(value.([]byte))
		} else {
			io.WriteString(w, "order not found")
		}
	}
}
