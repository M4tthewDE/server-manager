package internal

import (
	"net/http"

	"github.com/m4tthewde/server-manager/internal/routes"
)

func Run() error {
	http.HandleFunc("GET /{$}", routes.Root)
	http.HandleFunc("GET /status", routes.Stat)

	http.HandleFunc("GET /docker/{id}", routes.Details)
	http.HandleFunc("GET /docker/{id}/containerDetails", routes.ContainerDetails)
	http.HandleFunc("GET /docker/new", routes.ContainerForm)

	http.HandleFunc("POST /docker/new", routes.ContainerNew)
	http.HandleFunc("POST /docker/{id}/start", routes.ContainerStart)
	http.HandleFunc("POST /docker/{id}/stop", routes.ContainerStop)
	http.HandleFunc("POST /docker/{id}/remove", routes.ContainerRemove)
	http.HandleFunc("/docker/{id}/logs", routes.ContainerLogs)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	return http.ListenAndServe(":8080", nil)
}
