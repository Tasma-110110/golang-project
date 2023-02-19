package main

import (
	"html/template"
	"net/http"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("./static/*.html"))
}

func main() {
	http.HandleFunc("/", IndexHandler)
	http.ListenAndServe(":7777", nil)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.html", nil)
}
