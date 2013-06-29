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

	makeShield(w)
}

func index(w http.ResponseWriter, r *http.Request) {
	const idx = `
<html>
<head><title>Buckler</title></head>
<body>
<img src="/v1">
</body>
</html>
`
	fmt.Fprintf(w, idx)
}

func main() {
	ip := os.Getenv("OPENSHIFT_INTERNAL_IP")
	http.HandleFunc("/v1", buckle)
	http.HandleFunc("/", index)
	log.Println("Listening on port 8080");
	http.ListenAndServe(fmt.Sprintf("%s:8080", ip), nil)
}
