package main

import (
	"log"
	"os"
	"text/template"
)

func main() {
	t := template.New("main.tmpl")
	t = template.Must(t.ParseGlob("templates/*.tmpl"))
	err := t.Execute(os.Stdout, nil)
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}
}
