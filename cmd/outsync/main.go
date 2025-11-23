package main

import (
	"fmt"
	"outsync/internal/database"
	"context"
)

func main() {
	fmt.Println("Starting OutSync ...!")
	ctx := context.Background()
	
	err := database.ApplySchema(ctx)
	if err != nil {
		fmt.Println("Error applying schema:", err)
		return
	}
	conn ,err := database.Connect(ctx)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer conn.Close(ctx)	
	err = database.CreateUserWithEvent(ctx, conn, `{"prompt":"build a yc backup company"}`)
	if err != nil {
		fmt.Println("Error creating user:", err)
		return
	}
	fmt.Println("User with prompt created successfully")

}

