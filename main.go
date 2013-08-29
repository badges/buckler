package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"./shield"
	"github.com/droundy/goopt"
)

var (
	wsReplacer    = strings.NewReplacer("__", "_", "_", " ")
	revWsReplacer = strings.NewReplacer(" ", "_", "_", "__", "-", "--")

	// set last modifed to server startup. close enough to release.
	lastModified    = time.Now()
	lastModifiedStr = lastModified.UTC().Format(http.TimeFormat)
	oneYear         = time.Duration(8700) * time.Hour

	staticPath = "static"
)

func shift(s []string) ([]string, string) {
	return s[1:], s[0]
}

func invalidRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("bad request", r.URL.String())
	http.Error(w, "bad request", 400)
}

func parseFileName(name string) (d shield.Data, err error) {
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
	c, err := shield.GetColor(cp)
	if err != nil {
		return
	}

	d = shield.Data{newParts[0], newParts[1], c}
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

	shield.PNG(w, d)
}

const basePkg = "github.com/gittip/img.shields.io"

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(staticPath, "index.html"))
}

func favicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(staticPath, "favicon.png"))
}

func fatal(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func cliMode(vendor string, status string, color string, args []string) {
	// if any of vendor, status or color is given, all must be
	if (vendor != "" || status != "" || color != "") &&
		!(vendor != "" && status != "" && color != "") {
		fatal("You must specify all of vendor, status, and color")
	}

	if vendor != "" {
		c, err := shield.GetColor(color)
		if err != nil {
			fatal("Invalid color: " + color)
		}
		d := shield.Data{vendor, status, c}

		name := fmt.Sprintf("%s-%s-%s.png", revWsReplacer.Replace(vendor),
			revWsReplacer.Replace(status), color)

		if len(args) > 1 {
			fatal("You can only specify one output file name")
		}

		if len(args) == 1 {
			name = args[0]
		}

		// default to standard out
		f := os.Stdout
		if name != "-" {
			f, err = os.Create(name)
			if err != nil {
				fatal("Unable to create file: " + name)
			}
		}

		shield.PNG(f, d)
		return
	}

	// generate based on command line file names
	for i := range args {
		name := args[i]
		d, err := parseFileName(name)
		if err != nil {
			fatal(err.Error())
		}

		f, err := os.Create(name)
		if err != nil {
			fatal(err.Error())
		}
		shield.PNG(f, d)
	}
}

func usage() string {
	u := `Usage: %s [-h HOST] [-p PORT]
       %s [-v VENDOR -s STATUS -c COLOR] <FILENAME>

%s`
	return fmt.Sprintf(u, os.Args[0], os.Args[0], goopt.Help())
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

	goopt.Usage = usage

	// common options
	dataDir := goopt.String([]string{"-d", "--data-dir"}, "data", "data dir containing base PNG files and font")

	// server mode options
	host := goopt.String([]string{"-h", "--host"}, hostEnv, "host ip address to bind to")
	port := goopt.Int([]string{"-p", "--port"}, p, "port to listen on")

	// cli mode
	vendor := goopt.String([]string{"-v", "--vendor"}, "", "vendor for cli generation")
	status := goopt.String([]string{"-s", "--status"}, "", "status for cli generation")
	color := goopt.String([]string{"-c", "--color", "--colour"}, "", "color for cli generation")
	goopt.Parse(nil)

	args := goopt.Args

	shield.Init(*dataDir)

	// if any of the cli args are given, or positional args remain, assume cli
	// mode.
	if len(args) > 0 || *vendor != "" || *status != "" || *color != "" {
		cliMode(*vendor, *status, *color, args)
		return
	}
	// normalize for http serving
	if *host == "*" {
		*host = ""
	}

	http.HandleFunc("/v1/", buckle)
	http.HandleFunc("/favicon.png", favicon)
	http.HandleFunc("/", index)

	log.Println("Listening on port", *port)
	http.ListenAndServe(*host+":"+strconv.Itoa(*port), nil)
}
