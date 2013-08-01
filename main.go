package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	wsReplacer = strings.NewReplacer("__", "_", "_", " ")
)

func shift(s []string) ([]string, string) {
	return s[1:], s[0]
}

func invalidRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("bad request", r.URL.String())
	http.Error(w, "bad request", 400)
}

func buckle(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) != 3 {
		invalidRequest(w, r)
		return
	}

	imageName := wsReplacer.Replace(parts[2])
	imageParts := strings.Split(imageName, "-")

	newParts := []string{}
	for len(imageParts) > 0 {
		var head, right string
		imageParts, head = shift(imageParts)

		// if starts with - append to previous
		if len(head) == 0 && len(newParts) > 0 {
			left := ""
			if len(newParts) > 0 {
				left = newParts[len(newParts)-1]
				newParts = newParts[:len(newParts)-1]
			}

			// trailing -- is going to break color anyways so don't worry
			imageParts, right = shift(imageParts)

			head = strings.Join([]string{left, right}, "-")
		}

		newParts = append(newParts, head)
	}

	if len(newParts) != 3 {
		invalidRequest(w, r)
		return
	}

	if !strings.HasSuffix(newParts[2], ".png") {
		invalidRequest(w, r)
		return
	}

	c, ok := Colors[newParts[2][0:len(newParts[2])-4]]
	if !ok {
		invalidRequest(w, r)
		return
	}

	t, err := time.Parse(time.RFC1123, r.Header.Get("if-modified-since"))
	if err == nil && !t.Before(lastModified) {
		w.WriteHeader(304)
		return
	}

	d := Data{newParts[0], newParts[1], c}
	makePngShield(w, d)
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func main() {
	ip := os.Getenv("HOST")
	port := os.Getenv("PORT")
	http.HandleFunc("/v1/", buckle)
	http.HandleFunc("/", index)
	log.Println("Listening on port", port)
	http.ListenAndServe(ip+":"+port, nil)
}
