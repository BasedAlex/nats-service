package db

import (
	"context"
	"encoding/json"

	"github.com/basedalex/nats-service/model"
	"github.com/jackc/pgx/v5"
)

type Postgres struct {
	db *pgx.Conn
}

func NewPostgres(ctx context.Context, dbConnect string) (*Postgres, error) {
	db, err := pgx.Connect(ctx, dbConnect)
	if err != nil {
		return nil, err
	}
	return &Postgres{db:db}, nil
}

func (db *Postgres) CreateOrder(ctx context.Context, order model.OrderData) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}
	stmt := "INSERT INTO orders (order_id, order_data) VALUES ($1, $2);"
	_, err = db.db.Exec(ctx, stmt, order.OrderUid, data)
	if err != nil {
		return err
	}
	return nil
}

func (db *Postgres) GetAllOrders(ctx context.Context) <-chan GetAllOrdersResult {
	const bucketSize = 100
	orderCh := make(chan GetAllOrdersResult)
	go func (){
		defer close(orderCh)
		var id string
		for {
			stmt := "SELECT * FROM orders WHERE order_id > $1 ORDER BY order_id LIMIT $2"
			rows, err := db.db.Query(ctx, stmt, id, bucketSize)
			if err != nil {
				orderCh <- GetAllOrdersResult{
					Err: err,
				}
				return
			}
			r := GetAllOrdersResult{}
			count := 0
			for rows.Next() {
				count++
				err = rows.Scan(&r.Data.ID, &r.Data.OrderData)
				if err != nil {
					orderCh <- GetAllOrdersResult{
						Err: err,
					}
					return
				}
				orderCh <- r
			}

			if err := rows.Err(); err != nil {
				orderCh <- GetAllOrdersResult{
					Err: err,
				}
				return
			}
			if count < bucketSize {
				return 
			}

		}
	}()

	return orderCh
}

type GetAllOrdersResult struct {
	Data model.DataItem
	Err error
}