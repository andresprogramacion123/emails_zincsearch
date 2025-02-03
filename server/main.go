package main

import (
	"fmt"
	"log"
	"net/http"
	"julian_zincsearch/server/controllers"
)

func main() {
	port := 8080

	app := controllers.App{}

	server := app.Routes()

	log.Printf("Server running on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), server))
}
