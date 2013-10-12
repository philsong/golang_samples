package main

import (
	"html/template"
	"net/http"
)

var hogeTmpl = template.Must(template.New("hoge").ParseFiles("base.html", "hoge.html"))

func hogeHandler(w http.ResponseWriter, r *http.Request) {
	hogeTmpl.ExecuteTemplate(w, "base", "Hoge")
}

var piyoTmpl = template.Must(template.New("piyo").ParseFiles("base.html", "piyo.html"))

func piyoHandler(w http.ResponseWriter, r *http.Request) {
	piyoTmpl.ExecuteTemplate(w, "base", "Piyo")
}

func main() {
	// hoge
	http.HandleFunc("/", hogeHandler)
	http.HandleFunc("/hoge", hogeHandler)

	// piyo
	http.HandleFunc("/piyo", piyoHandler)

	http.ListenAndServe(":8080", nil)
}
