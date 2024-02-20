package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		fmt.Fprintf(w, "Hello worwodowild!")
	})

	port := "3000"
	log.Printf("Listening on port %s", port)
	http.ListenAndServe(":"+port, mux)
}
