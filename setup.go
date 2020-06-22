package main

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

// Setup all connection including to nats server and route
func (s *Server) setup() error {

	var err error
	s.token = os.Getenv("APP_KEY")
	if s.token == "" {
		return errors.New("APP_KEY NOT SET")
	}

	s.maxImageSize, err = strconv.Atoi(os.Getenv("MAX_IMAGE_SIZE"))
	if err != nil {
		log.Println("Failed to get MAX_IMAGE_SIZE from .env using default 10M")
		s.maxImageSize = 10
	}

	// From lilliput readme
	s.imageExtSupport = map[string]bool{
		".jpeg": true,
		".webp": true,
		".png":  true,
		".jpg":  true,
	}

	s.router = mux.NewRouter()
	s.route()
	return nil
}
