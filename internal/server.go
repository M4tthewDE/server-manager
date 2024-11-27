package internal

import (
	"html/template"
	"log"
	"net/http"

	"github.com/m4tthewde/server-manager/internal/status"
)

func Run() error {
	http.HandleFunc("/{$}", root)
	http.HandleFunc("/status", status.Stat)
	http.HandleFunc("/favicon.ico", favicon)
	http.HandleFunc("/htmx.min.js", htmx)

	return http.ListenAndServe(":8080", nil)
}

var rootTemplate = template.Must(template.ParseFiles("static/root.html"))

func root(w http.ResponseWriter, r *http.Request) {
	err := rootTemplate.Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
}

func favicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/favicon.ico")
}

func htmx(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/htmx.min.js")
}
