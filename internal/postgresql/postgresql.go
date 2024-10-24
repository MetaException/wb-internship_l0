package postgresql

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Postgres struct {
	logger *logrus.Logger
	conn   *pgx.Conn
	config *Config
}

func New(logger *logrus.Logger, config *Config) *Postgres {

	conn, err := pgx.Connect(context.Background(), config.URL)
	if err != nil {
		logger.WithError(err).Fatal("unable to connect to db")
	}

	logger.Info("successfully connected to database")

	return &Postgres{
		logger: logger,
		conn:   conn,
		config: config,
	}
}

func (pg *Postgres) Close() {
	pg.conn.Close(context.Background())
}

func (pg *Postgres) AddToDb(uid string, json_data json.RawMessage) error {
	_, err := pg.conn.Exec(context.Background(), "insert into orders(order_uid, data) values($1, $2)", uid, string(json_data))

	if err != nil {
		pg.logger.WithError(err).Error("failed to add to db")
		return errors.WithStack(err)
	}

	pg.logger.Info("successfully added the row to db")

	return nil
}

func (pg *Postgres) GetAllData() (map[string]interface{}, error) {
	rows, err := pg.conn.Query(context.Background(), "select * from orders")

	if err != nil {
		pg.logger.WithError(err).Error("get all data failed")
		return nil, errors.WithStack(err)
	}

	pg.logger.Info("successfully received all data from db")

	orders := make(map[string]interface{})
	for rows.Next() {

		var uid string
		var json []byte

		err := rows.Scan(&uid, &json)
		if err != nil {
			pg.logger.WithError(err).Error("failed to scan data from row")
			return nil, errors.WithStack(err)
		}

		orders[uid] = json
	}

	return orders, nil
}
