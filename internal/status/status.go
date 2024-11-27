package status

import (
	"html/template"
	"log"
	"net/http"
)

var templ = template.Must(template.ParseFiles("static/status.html"))

type Status struct {
	Memory Memory
	Docker Docker
}

func NewStatus() (*Status, error) {
	memory, err := FetchMemory()
	if err != nil {
		return nil, err
	}

	docker, err := FetchDocker()
	if err != nil {
		return nil, err
	}

	return &Status{Memory: *memory, Docker: *docker}, nil
}

func Stat(w http.ResponseWriter, r *http.Request) {
	status, err := NewStatus()
	if err != nil {
		log.Println(err)
	}

	err = templ.Execute(w, status)
	if err != nil {
		log.Println(err)
	}
}
