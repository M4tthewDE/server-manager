package success

import (
	"html/template"
	"log"
	"net/http"
)

var successTempl = template.Must(template.ParseFiles("static/success.html"))

func ShowSuccess(w http.ResponseWriter, msg string) {
	err := successTempl.Execute(w, msg)
	if err != nil {
		log.Println(err)
	}
}
