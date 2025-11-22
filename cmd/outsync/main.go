package main

import (
	"fmt"
	"outsync/internal/database"
	"outsync/internal/config"
	"context"
)

func main() {
	fmt.Println("Starting OutSync ...!")
	config = config.LoadConfig()
	database.NewConnection(context.Background())
	
}

