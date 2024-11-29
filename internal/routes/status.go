package routes

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"github.com/m4tthewde/server-manager/internal/docker"
	"github.com/m4tthewde/server-manager/internal/memory"
)

var templ = template.Must(template.ParseFiles("static/status.html"))

type Status struct {
	Memory memory.Memory
	Docker docker.Docker
}

func NewStatus(ctx context.Context) (*Status, error) {
	memory, err := memory.FetchMemory()
	if err != nil {
		return nil, err
	}

	docker, err := docker.FetchDocker(ctx)
	if err != nil {
		return nil, err
	}

	return &Status{Memory: *memory, Docker: *docker}, nil
}

func Stat(w http.ResponseWriter, r *http.Request) {
	status, err := NewStatus(r.Context())
	if err != nil {
		log.Println(err)
	}

	err = templ.Execute(w, status)
	if err != nil {
		log.Println(err)
	}
}
