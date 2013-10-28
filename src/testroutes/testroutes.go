package main

import (
	"fmt"
	"github.com/drone/routes"
	"net/http"
	"os"
)

func Whoami(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	lastName := params.Get(":last")
	firstName := params.Get(":first")
	fmt.Fprintf(w, "you are %s %s", firstName, lastName)
}

func main() {
	mux := routes.New()
	mux.Get("/:last/:first", Whoami)
	pwd, _ := os.Getwd()
	mux.Static("/static", pwd)
	http.Handle("/", mux)
	http.ListenAndServe(":8088", nil)
}
