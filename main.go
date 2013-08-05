package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	wsReplacer = strings.NewReplacer("__", "_", "_", " ")

	// set last modifed to server startup. close enough to release.
	lastModified    = time.Now()
	lastModifiedStr = lastModified.UTC().Format(http.TimeFormat)
	oneYear         = time.Duration(8700) * time.Hour
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

	cp := newParts[2][0 : len(newParts[2])-4]
	c, ok := Colors[cp]
	if !ok {
		c, ok = hexColor(cp)
		if !ok {
			invalidRequest(w, r)
			return
		}
	}

	t, err := time.Parse(time.RFC1123, r.Header.Get("if-modified-since"))
	if err == nil && lastModified.Before(t.Add(1*time.Second)) {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	w.Header().Add("content-type", "image/png")
	w.Header().Add("expires", time.Now().Add(oneYear).Format(time.RFC1123))
	w.Header().Add("cache-control", "public")
	w.Header().Add("last-modified", lastModifiedStr)

	d := Data{newParts[0], newParts[1], c}
	makePngShield(w, d)
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func favicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/favicon.png")
}

func main() {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	hostArg := flag.String("host", "*", "host ip address to bind to")
	portArg := flag.String("port", "8080", "port to listen on")
	flag.Parse()

	hostSet := false
	portSet := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "host" {
			hostSet = true
		}

		if f.Name == "port" {
			portSet = true
		}
	})

	if hostSet || host == "" {
		host = *hostArg
	}

	if portSet || port == "" {
		port = *portArg
	}

	if host == "*" {
		host = ""
	}

	http.HandleFunc("/v1/", buckle)
	http.HandleFunc("/favicon.png", favicon)
	http.HandleFunc("/", index)

	log.Println("Listening on port", port)
	http.ListenAndServe(host+":"+port, nil)
}
