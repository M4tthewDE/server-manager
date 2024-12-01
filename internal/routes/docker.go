package routes

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/m4tthewde/server-manager/internal/docker"
	"github.com/m4tthewde/server-manager/internal/error"
	"github.com/m4tthewde/server-manager/internal/success"
)

var detailsTempl = template.Must(template.ParseFiles("static/docker/details.html"))
var containerDetailsTempl = template.Must(template.ParseFiles("static/docker/containerDetails.html"))
var newTempl = template.Must(template.ParseFiles("static/docker/new.html"))

func Details(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := detailsTempl.Execute(w, id)
	if err != nil {
		log.Println(err)
	}
}

func ContainerForm(w http.ResponseWriter, r *http.Request) {
	err := newTempl.Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
}

func ContainerDetails(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	container, err := docker.FindContainer(r.Context(), id)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", 302)
	}

	err = containerDetailsTempl.Execute(w, container)
	if err != nil {
		log.Println(err)
	}
}

func ContainerNew(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		error.ShowError(w, err)
		return
	}

	image := r.FormValue("image")
	if image == "" {
		error.ShowError(w, errors.New("no image provided"))
		return
	}

	version := r.FormValue("version")
	if version == "" {
		error.ShowError(w, errors.New("no version provided"))
		return
	}

	id, err := docker.CreateContainer(r.Context(), image, version)
	if err != nil {
		error.ShowError(w, err)
		return
	}

	w.Header().Add("HX-Location", "/docker/"+id)
}

func ContainerStart(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := docker.StartContainer(r.Context(), id)
	if err != nil {
		error.ShowError(w, err)
		return
	}

	success.ShowSuccess(w, "")
}

func ContainerStop(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := docker.StopContainer(r.Context(), id)
	if err != nil {
		error.ShowError(w, err)
		return
	}

	success.ShowSuccess(w, "")
}

func ContainerRemove(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := docker.RemoveContainer(r.Context(), id)
	if err != nil {
		error.ShowError(w, err)
		return
	}

	w.Header().Add("HX-Location", "/")
}

func ContainerLogs(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	clientGone := r.Context().Done()
	rc := http.NewResponseController(w)

	logsChannel := make(chan docker.LogMessage)
	go docker.StreamLogs(r.Context(), id, logsChannel)

	for {
		select {
		case <-clientGone:
			return
		case msg := <-logsChannel:
			if msg.Error != nil {
				_, err := fmt.Fprintf(w, "event: ErrorEvent\ndata: <span>Error: %s</span>\n\n", msg.Error.Error())
				if err != nil {
					log.Println(err)
					return
				}
			} else {
				_, err := fmt.Fprintf(w, "event: LogEvent\ndata: <br><span>%s</span>\n\n", msg.Text)
				if err != nil {
					log.Println(err)
					return
				}
			}

			err := rc.Flush()
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}
