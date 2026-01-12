package main

import (
	"flag"
	"fmt"
	"os"
	"pixgrid/server"
)

func main() {
	port := flag.Int("port", 8080, "Port to run the server on")
	flag.Parse()

	srv := server.New()
	fmt.Printf("Starting pixgrid server on port %d...\n", *port)
	if err := srv.Start(*port); err != nil {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}
