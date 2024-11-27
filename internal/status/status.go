package status

import (
	"html/template"
	"log"
	"net/http"
)

var templ = template.Must(template.ParseFiles("static/status.html"))

func Stat(w http.ResponseWriter, r *http.Request) {
	memory, err := FetchMemory()
	if err != nil {
		log.Println(err)
	}

	err = templ.Execute(w, memory)
	if err != nil {
		log.Println(err)
	}
}
