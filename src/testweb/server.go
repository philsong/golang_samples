package main

import "github.com/codegangsta/martini"
import (
	"github.com/codegangsta/martini-contrib/auth"
)

func main() {
	m := martini.Classic()
	m.Get("/", func() string {
		return "Hello world!"
	})
	m.Use(auth.Basic("username", "secretpassword"))
	m.Use(martini.Static("./"))
	m.Run()
}
