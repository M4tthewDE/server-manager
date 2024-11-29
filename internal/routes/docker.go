package routes

import (
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
		log.Println(err)
		http.Redirect(w, r, "/", 302)
	}

	image := r.FormValue("image")
	if image == "" {
		http.Redirect(w, r, "/", 302)
	}

	version := r.FormValue("version")
	if image == "" {
		http.Redirect(w, r, "/", 302)
	}

	_, err = docker.CreateContainer(r.Context(), image, version)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", 302)
	}

	// TODO: redirect to details page of container
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
