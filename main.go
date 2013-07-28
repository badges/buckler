package main

import (
	"log"
	"net/http"
	"os"
)

func buckle(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	log.Printf("Requsted: %v", q)

	// arg validation goes here
	d := Data{q["v"][0], q["s"][0], q["c"][0]}
	makePngShield(w, d)
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func main() {
	ip := os.Getenv("HOST")
	port := os.Getenv("PORT")
	http.HandleFunc("/v1", buckle)
	http.HandleFunc("/", index)
	log.Println("Listening on port", port)
	http.ListenAndServe(ip+":"+port, nil)
}
