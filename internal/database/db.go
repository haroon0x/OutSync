package database

import (
	"fmt"
	"github.com/jackc/pgx/v5"
	"context"
	"os"
	"outsync/internal/config"
)

func Connect(ctx context.Context) (*pgx.Conn, error) {
	cfg := config.LoadConfig()
	conn,err := pgx.Connect(ctx , cfg.DatabaseUrl )
	if err != nil {
		fmt.Fprintf(os.Stderr,"unable to connect to databse: %v\n" , err)
		return nil, err
	}
	fmt.Println("Connected to database")
	return conn, err
}

func ApplySchema(ctx context.Context) error {
	conn,err := Connect(ctx)
	if err != nil {
		return err
	}
	schema , err := os.ReadFile("internal/database/schema.sql")
	if err != nil {
		return err
	}
	_,err = conn.Exec(ctx, string(schema))
	if err != nil {
		return err
	}
	fmt.Println("Database updated successfully")
	return err
}
