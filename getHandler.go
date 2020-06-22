package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func (s *Server) getHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		vars := mux.Vars(r)
		inputBuf, err := ioutil.ReadFile("uploads/" + vars["fileName"])
		if err != nil {
			log.Printf("failed to read input file, %s\n", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// change extension to match response content/type
		fileExt := filepath.Ext(vars["fileName"])
		if fileExt == ".jpg" {
			fileExt = ".jpeg"
		}

		// check if image need to be resized
		urlQuery := r.URL.Query()
		if len(urlQuery["h"]) > 0 || len(urlQuery["w"]) > 0 {
			width, err := strconv.Atoi(urlQuery["w"][0])
			if err != nil {
				width = 0
			}
			height, err := strconv.Atoi(urlQuery["h"][0])
			if err != nil {
				height = 0
			}
			log.Println(height, width)
			inputBuf, err = resize(inputBuf, height, width)
			if err != nil {
				log.Printf("failed to resize image, %s\n", err)
				http.Error(w, "failed to resize image", http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "image/"+strings.Replace(fileExt, ".", "", -1))
		w.Header().Set("Content-Length", strconv.Itoa(len(inputBuf)))
		if _, err := w.Write(inputBuf); err != nil {
			log.Println("unable to write image.")
		}

	}
}
