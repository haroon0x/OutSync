package database

import (
	"fmt"
	"github.com/jackc/pgx/v5"
	"context"
	"os"
	"outsync/internal/config"
)

func NewConnection(ctx context.Context) (*pgx.Conn, error) {
	cfg := config.LoadConfig()
	conn,err := pgx.Connect(ctx , cfg.DatabaseUrl )
	if err != nil {
		fmt.Fprintf(os.Stderr,"unable to connect to databse: %v\n" , err)
		return nil, err
	}
	return conn, err
	
}