package error

import (
	"html/template"
	"log"
	"net/http"
)

var errorTempl = template.Must(template.ParseFiles("static/error.html"))

func ShowError(w http.ResponseWriter, e error) {
	err := errorTempl.Execute(w, e.Error())
	if err != nil {
		log.Println(err)
	}
}
