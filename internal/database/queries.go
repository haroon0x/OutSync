package database

import (
	"context"
	"github.com/jackc/pgx/v5"
)

func CreateUserWithEvent (ctx context.Context , conn *pgx.Conn, data string ) (error) {
	tx, err := conn.Begin(ctx)
	if err!=nil{
		return err
	}
	defer tx.Rollback(ctx)
	var userId string
	err = tx.QueryRow(ctx,"INSERT INTO users (data) VALUES ($1) RETURNING id",data).Scan(&userId)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx , "INSERT INTO outbox_events (aggregate_type,aggregate_id,payload,status) VALUES ($1,$2,$3,$4)", "prompt", userId, data, "pending")
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}