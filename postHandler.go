package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

func (s *Server) postHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		err = r.ParseMultipartForm(10 << s.maxImageSize)
		if err != nil {
			log.Println("Parse multipartform ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		file, handler, err := r.FormFile("img")
		if err != nil {
			log.Println("Error Retrieving the File, err ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		log.Printf("Uploaded File: %+v File Size: %+v \n", handler.Filename, handler.Size)

		fileExt := filepath.Ext(handler.Filename)
		if _, ok := s.imageExtSupport[fileExt]; !ok {
			log.Println("File Extension Not supported : ", fileExt)
			http.Error(w, "File Extension Not supported", http.StatusBadRequest)
			return
		}

		// Create a temporary file within our temp-images directory that follows
		// a particular naming pattern
		tempFile, err := ioutil.TempFile("uploads", "*"+fileExt)
		if err != nil {
			log.Println("Tempfile cannot created ", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer tempFile.Close()

		// read all of the contents of our uploaded file into a
		// byte array
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println("Cannot Read File")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		// write this byte array to our temporary file
		_, err = tempFile.Write(fileBytes)
		if err != nil {
			log.Println("Cannot Write file ", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		// return that we have successfully uploaded our file!

		urlPath := "/image/" + tempFile.Name()
		w.Header().Set("Content-Type", "application/json")

		// Return fixed response change this if nescessary
		_, err = w.Write([]byte(fmt.Sprintf(`{"success":true,"urlPath": %s}`, urlPath)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}
