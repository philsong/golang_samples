package main

import (
	"log"
	"os"
	"text/template"
)

func main() {
	t := template.Must(template.ParseFiles("templates/main.tmpl", "templates/header.tmpl", "templates/footer.tmpl"))
	err := t.Execute(os.Stdout, nil)
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}
}
