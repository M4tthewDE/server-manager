package template

import "html/template"

var RootTemplate = template.Must(template.ParseFiles("static/root.html"))
var MemoryTemplate = template.Must(template.ParseFiles("static/memory.html"))
