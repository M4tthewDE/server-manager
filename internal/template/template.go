package template

import "html/template"

var RootTemplate = template.Must(template.ParseFiles("static/root.html"))
