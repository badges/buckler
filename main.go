package main

import (
	"os"
	"fmt"
	"log"
	"net/http"
)

func buckle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
	q := r.URL.Query()
	log.Printf("Requsted: %v", q)
}

func main() {
	ip := os.Getenv("OPENSHIFT_INTERNAL_IP")
	http.HandleFunc("/", buckle)
	log.Println("Listening on port 8080");
	http.ListenAndServe(fmt.Sprintf("%s:8080", ip), nil)
}
