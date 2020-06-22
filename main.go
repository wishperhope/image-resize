package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

// Server Struct Contain DB and Nats and token
// See : https://medium.com/@matryer/how-i-write-go-http-services-after-seven-years-37c208122831
type Server struct {
	router          *mux.Router
	token           string
	maxImageSize    int
	imageExtSupport map[string]bool
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error Loading .env File")
	}

	s := Server{}
	err = s.setup()
	if err != nil {
		panic(err)
	}

	// Setup HTTP Server
	httpPort := os.Getenv("HTTP_PORT")
	// Allow cors
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            true,
	},
	).Handler(s.router)
	log.Printf("REST API Started On %s\n", httpPort)
	log.Fatal(http.ListenAndServe(":"+httpPort, handler))

}
