package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"animai/templates"

	"github.com/a-h/templ"
)

func run(port int16, staticFolder *string) {
	log.SetPrefix("[SERVER] ")
	indexPage := templates.Index()

	http.Handle("/", templ.Handler(indexPage))
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir(*staticFolder))))

	log.Println("Listening on :", port)

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func main() {
	static := flag.String("static dir", "static/", "Static Folder")
	run(8080, static)
}
