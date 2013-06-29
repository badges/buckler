package main


import (
	"text/template"
	"net/http"
)

func makeShield (w http.ResponseWriter) {
	type Data struct {
		Vendor string
		Status string
		Color string
	}

	w.Header().Add("content-type", "image/svg+xml")
	d := Data{"Goes To", "11", "green"}

	t := template.Must(template.New("svg").ParseFiles("tmpl/shield.svg"))
	t.ExecuteTemplate(w, "svg", d)
}
