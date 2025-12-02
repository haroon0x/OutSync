package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type OutboxEvent struct {
	Id      string
	Payload string
	Status  string
}

func CreateUserWithEvent(ctx context.Context, conn *pgx.Conn, data string) error {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var userId string
	err = tx.QueryRow(ctx, "INSERT INTO users (data) VALUES ($1) RETURNING id", data).Scan(&userId)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, "INSERT INTO outbox_events (aggregate_type,aggregate_id,payload,status) VALUES ($1,$2,$3,$4)", "prompt", userId, data, "pending")
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func GetPendingEvents(ctx context.Context, conn *pgx.Conn) ([]OutboxEvent, error) {
	res, err := conn.Query(ctx, "SELECT id , payload , status FROM outbox_events WHERE status='pending' LIMIT 10")
	if err != nil {
		return nil , err
	}
	defer res.Close()
	var events []OutboxEvent
	for res.Next() {
		var e OutboxEvent
		err := res.Scan(&e.Id, &e.Payload, &e.Status)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	if res.Err() != nil {
		return nil, res.Err()
	}
	return events, nil
}

func MarkEventProcessed(ctx context.Context, conn *pgx.Conn, eventID string) error {
	_, err := conn.Exec(ctx, "UPDATE outbox_events SET status='processed' where id = $1", eventID)
	if err != nil {
		return err
	}
	return nil
}
