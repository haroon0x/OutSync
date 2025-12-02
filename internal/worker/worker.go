package worker

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/jackc/pgx/v5"
	"outsync/internal/database"
	"fmt"
	"time"
)


func Poll(ctx context.Context , conn *pgx.Conn , rdb *redis.Client){
	fmt.Println("worker started Polling for events ...!")
	
	for {
		events , err := database.GetPendingEvents(ctx , conn)
		if err != nil {
			fmt.Println("Error polling for events:", err)
			time.Sleep(time.Second * 5)
			continue
		}
		if len(events) == 0 {
			fmt.Println("No events found")
			time.Sleep(time.Second * 5)
			continue
		}
		fmt.Println("Found" , len(events) , "events" , (events)) 
		for _,event := range events {
			fmt.Printf("Processing event %s\n", event.Id)
			listlen , err := rdb.RPush(ctx,"events_queue" , event.Payload).Result()
			if err != nil {
				fmt.Println("Error pushing event to queue:", err)
				continue
			}
			fmt.Println("Pushed event to queue. List length:", listlen)
			err = database.MarkEventProcessed(ctx , conn, event.Id)
			if err != nil {
				fmt.Println("Error marking event as processed:", err)
				continue
			}
		}
 		
	}

}