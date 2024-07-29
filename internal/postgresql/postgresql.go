package postgresql

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type Postgres struct {
	logger *logrus.Logger //TODO
	conn   *pgx.Conn
}

func New() *Postgres {

	conn, err := pgx.Connect(context.Background(), "postgresql://localdbuser:localdbuserpass@localhost:5433/wbl0")
	if err != nil {
		logrus.Fatal("unable to connect to db: %w", err)
	}

	return &Postgres{
		logger: logrus.New(),
		conn:   conn,
	}
}

func (pg *Postgres) Close() {
	pg.conn.Close(context.Background())
}

func (pg *Postgres) AddToDb(uid string, json_data json.RawMessage) error {
	_, err := pg.conn.Exec(context.Background(), "insert into orders(order_uid, data) values($1, $2)", uid, string(json_data))

	if err != nil {
		return fmt.Errorf("failed to add to db: %w", err)
	}
	return nil
}

func (pg *Postgres) GetAllData() (map[string]interface{}, error) {
	rows, err := pg.conn.Query(context.Background(), "select * from orders")

	if err != nil {
		return nil, fmt.Errorf("queryRow failed: %w", err)
	}

	orders := make(map[string]interface{})
	for rows.Next() {

		var uid string
		var json []byte

		err := rows.Scan(&uid, &json)
		if err != nil {
			return nil, fmt.Errorf("failed to scan data from row: %w", err)
		}

		orders[uid] = json
	}

	return orders, nil
}
