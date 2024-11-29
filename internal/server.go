package internal

import (
	"html/template"
	"log"
	"net/http"

	"github.com/m4tthewde/server-manager/internal/routes"
)

func Run() error {
	http.HandleFunc("GET /{$}", root)
	http.HandleFunc("GET /status", routes.Stat)
	http.HandleFunc("GET /docker/{id}", routes.Details)
	http.HandleFunc("GET /docker/{id}/containerDetails", routes.ContainerDetails)
	http.HandleFunc("GET /docker/new", routes.ContainerForm)
	http.HandleFunc("POST /docker/new", routes.ContainerNew)
	http.HandleFunc("POST /docker/{id}/start", routes.ContainerStart)
	http.HandleFunc("POST /docker/{id}/stop", routes.ContainerStop)
	http.HandleFunc("POST /docker/{id}/remove", routes.ContainerRemove)

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
