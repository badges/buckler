package main

import (
	"net/http"
	"text/template"
)

type Data struct {
	Vendor string
	Status string
	Color  string
}

func makeShield(w http.ResponseWriter, d Data) {
	w.Header().Add("content-type", "image/svg+xml")

	t := template.Must(template.New("svg").ParseFiles("tmpl/shield.svg"))
	t.ExecuteTemplate(w, "svg", d)
}
