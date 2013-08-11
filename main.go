package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/droundy/goopt"
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

func parseFileName(name string) (d Data, err error) {
	imageName := wsReplacer.Replace(name)
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
		err = errors.New("Invalid file name")
		return
	}

	if !strings.HasSuffix(newParts[2], ".png") {
		err = errors.New("Unknown file type")
		return
	}

	cp := newParts[2][0 : len(newParts[2])-4]
	c, err := getColor(cp)
	if err != nil {
		return
	}

	d = Data{newParts[0], newParts[1], c}
	return
}

func buckle(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) != 3 {
		invalidRequest(w, r)
		return
	}

	d, err := parseFileName(parts[2])
	if err != nil {
		invalidRequest(w, r)
		return
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

	makePngShield(w, d)
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func favicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/favicon.png")
}

func main() {
	hostEnv := os.Getenv("HOST")
	portEnv := os.Getenv("PORT")

	// default to environment variable values (changes the help string :( )
	if hostEnv == "" {
		hostEnv = "*"
	}

	p := 8080
	if portEnv != "" {
		p, _ = strconv.Atoi(portEnv)
	}

	// server mode options
	host := goopt.String([]string{"-h", "--host"}, hostEnv, "host ip address to bind to")
	port := goopt.Int([]string{"-p", "--port"}, p, "port to listen on")

	// cli mode
	vendor := goopt.String([]string{"-v", "--vendor"}, "", "vendor for cli generation")
	status := goopt.String([]string{"-s", "--status"}, "", "status for cli generation")
	color := goopt.String([]string{"-c", "--color", "--colour"}, "", "color for cli generation")
	goopt.Parse(nil)

	if *host == "*" {
		*host = ""
	}

	args := goopt.Args

	if *vendor != "" {
		c, err := getColor(*color)
		if err != nil {
			log.Fatal(err)
		}
		d := Data{*vendor, *status, c}

		// XXX could escape here
		name := *vendor + "-" + *status + "-" + *color + ".png"

		if len(args) > 1 {
			log.Fatal("You can only specify one output file name")
		}

		if len(args) == 1 {
			name = args[0]
		}

		// default to standard out
		f := os.Stdout
		if name != "-" {
			f, err = os.Create(name)
			if err != nil {
				log.Fatal(err)
			}
		}

		makePngShield(f, d)
		return
	}

	// command line image generation
	if len(args) > 0 {
		for i := range args {
			name := args[i]
			d, err := parseFileName(name)
			if err != nil {
				log.Fatal(err)
			}

			f, err := os.Create(name)
			if err != nil {
				log.Fatal(err)
			}
			makePngShield(f, d)
		}
		return
	}

	http.HandleFunc("/v1/", buckle)
	http.HandleFunc("/favicon.png", favicon)
	http.HandleFunc("/", index)

	log.Println("Listening on port", *port)
	http.ListenAndServe(*host+":"+strconv.Itoa(*port), nil)
}
