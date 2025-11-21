package main

import (

	"fmt"
	"outsync/internal/config"
)

func main() {
	fmt.Println("Starting OutSync ...!")
	config.LoadConfig()
	
}

