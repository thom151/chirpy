package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Hello World")

	serverMux := http.NewServeMux()

	//SERVER
	httpServer := http.Server{}

	httpServer.Handler = serverMux
	httpServer.Addr = ":8080"

	err := httpServer.ListenAndServe()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
