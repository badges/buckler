package main

import (
	"os"
	"fmt"
	"log"
	"net/http"
)

func buckle(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	log.Printf("Requsted: %v", q)

	// arg validation goes here
	d := Data{q["v"][0], q["s"][0], q["c"][0]}
	makeShield(w, d)
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func main() {
	ip := os.Getenv("OPENSHIFT_INTERNAL_IP")
	http.HandleFunc("/v1", buckle)
	http.HandleFunc("/", index)
	log.Println("Listening on port 8080");
	http.ListenAndServe(fmt.Sprintf("%s:8080", ip), nil)
}
