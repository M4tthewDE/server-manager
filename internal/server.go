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

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	return http.ListenAndServe(":8080", nil)
}

var rootTemplate = template.Must(template.ParseFiles("static/root.html"))

func root(w http.ResponseWriter, r *http.Request) {
	err := rootTemplate.Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
}
