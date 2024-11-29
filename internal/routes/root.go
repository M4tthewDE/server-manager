package routes

import (
	"html/template"
	"log"
	"net/http"
)

var rootTemplate = template.Must(template.ParseFiles("static/root.html"))

func Root(w http.ResponseWriter, r *http.Request) {
	err := rootTemplate.Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
}
