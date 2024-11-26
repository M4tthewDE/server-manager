package internal

import (
	"log"
	"net/http"

	"github.com/m4tthewde/server-manager/internal/memory"
	"github.com/m4tthewde/server-manager/internal/template"
)

func Run() error {
	http.HandleFunc("/{$}", root)
	http.HandleFunc("/memory", mem)
	http.HandleFunc("/favicon.ico", favicon)
	http.HandleFunc("/htmx.min.js", htmx)

	return http.ListenAndServe(":8080", nil)
}

func root(w http.ResponseWriter, r *http.Request) {
	err := template.RootTemplate.Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
}

func mem(w http.ResponseWriter, r *http.Request) {
	memory, err := memory.FetchMemory()
	if err != nil {
		log.Println(err)
	}

	err = template.MemoryTemplate.Execute(w, memory)
	if err != nil {
		log.Println(err)
	}
}

func favicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/favicon.ico")
}

func htmx(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/htmx.min.js")
}
