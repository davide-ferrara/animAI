package server

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/davide-ferrara/animAI/templates/compiled"
)

func ServerRun() {
	static := flag.String("d", "../static", "Static Folder")
	component := compiled.Hello("Dave")

	http.Handle("/", templ.Handler(component))
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir(*static))))

	fmt.Println("Listening on :8080")
	http.ListenAndServe(":8080", nil)
}

func main() {
	ServerRun()
}
