package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
)

func customFileServer(fs http.FileSystem) http.Handler {
	fileServer := http.FileServer(fs)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fs.Open(path.Clean(r.URL.Path)) // Do not allow path traversals.
		if os.IsNotExist(err) {
			http.ServeFile(w, r, "./static/index.html")
			return
		}
		fileServer.ServeHTTP(w, r)
	})
}

func listFilesFromServer(w http.ResponseWriter, r *http.Request) {
	galleryFiles, err := os.ReadDir("./static/gallery")
	sliderFiles, sliderErr := os.ReadDir("./static/slider")
	if err != nil || sliderErr != nil {
		log.Fatal(err)
	}
	response := make(map[string][]string)

	for _, file := range galleryFiles {
		response["galleryPhotos"] = append(response["galleryPhotos"], fmt.Sprintf("https://cdn.sjpc.me/gallery/%s", file.Name()))
	}
	for _, file := range sliderFiles {
		response["sliderPhotos"] = append(response["sliderPhotos"], fmt.Sprintf("https://cdn.sjpc.me/slider/%s", file.Name()))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}

func main() {
	fs := customFileServer(http.Dir("./static"))
	http.HandleFunc("/list", listFilesFromServer)
	http.Handle("/", fs)
	fmt.Print("Listening on 4000")
	err := http.ListenAndServe(":4000", nil)
	if err != nil {
		log.Fatal(err)
	}

}
