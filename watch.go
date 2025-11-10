package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

func runServer() (*os.Process, error) {
	log.Println("WATCHER Building server binary")
	buildCmd := exec.Command("make", "server")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	if err := buildCmd.Run(); err != nil {
		log.Println("WATCHER ERROR executing make server:", err)
	} else {
		log.Println("WATCHER Web Server recompiled! Browser reload required.")
	}

	runServerCmd := exec.Command("./webserver")
	runServerCmd.Stdout = os.Stdout
	runServerCmd.Stderr = os.Stderr

	// Not blocking
	if err := runServerCmd.Start(); err != nil {
		log.Println("WATCHER: Failed to start server:", err)
		return nil, errors.New("err")
	}
	process := runServerCmd.Process
	log.Print("WATCHER Server running on PID: ", process.Pid)
	return process, nil
}

func stopServer(serverProcess *os.Process) {
	if serverProcess == nil {
		return
	}
	log.Printf("WATCHER: Stopping server process (PID: %d)...", serverProcess.Pid)
	if err := serverProcess.Signal(syscall.SIGTERM); err != nil {
		log.Println("WATCHER: Failed to stop server, killing:", err)
		if err := serverProcess.Kill(); err != nil {
			log.Println("WATCHER Could not kill the server process!")
		}
	}
	if _, err := serverProcess.Wait(); err != nil {
		log.Println("WATCHER Could not wait server process!")
	}
	serverProcess = nil
}

func Watch() {
	log.SetPrefix("[WATCHER] ")

	for {
		var serverProcess *os.Process
		var err error

		serverProcess, err = runServer()
		if err != nil {
			log.Fatal("WATCHER ERROOR Server could not Run!")
		}

		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		// defer watcher.Close()
		done := make(chan struct{})

		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}
					log.Print("New event: ", event.Name)

					if strings.HasSuffix(event.Name, "_templ.go") || strings.HasSuffix(event.Name, "bundle.js") {
						continue // Salta l'evento
					}

					// Esegui la compilazione solo per eventi di WRITE, CREATE, REMOVE, ecc.
					if event.Has(fsnotify.Write) {
						// 2. FILTRAGGIO ESECUZIONE
						if strings.HasSuffix(event.Name, ".js") || strings.HasSuffix(event.Name, ".css") {
							// A. MODIFICA ASSET (JS/CSS): Solo ricompilazione Esbuild
							log.Printf("Asset change detected in: %s. Recompiling JS/CSS...", event.Name)

							compileJs := exec.Command("make", "es")
							compileJs.Stdout = os.Stdout
							compileJs.Stderr = os.Stderr

							if err := compileJs.Run(); err != nil {
								log.Println("ERROR executing make es:", err)
							} else {
								log.Println("JS/CSS Recompiled! Browser reload required.")
							}

						} else if strings.HasSuffix(event.Name, ".go") || strings.HasSuffix(event.Name, ".templ") {
							compileTempl := exec.Command("make", "templ")
							compileTempl.Stdout = os.Stdout
							compileTempl.Stderr = os.Stderr

							if err := compileTempl.Run(); err != nil {
								log.Println("WATCHER ERROR executing make templ:", err)
								continue
							}
							log.Println("WATCHER Templ Recompiled!")

							close(done)
							return
						}
					}
				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}
					log.Println("Error:", err)
				}
			}
		}()

		directoryToWatch := []string{
			".",
			"server",
			"src",
			"templates",
		}

		for _, dir := range directoryToWatch {
			// Check is dir extists
			_, err := os.Stat(dir)
			if err != nil {
				log.Printf("WATCHER %s is not a valid dir and will be not watched", dir)
				continue
			}
			if err = watcher.Add(dir); err != nil {
				log.Fatal(err)
			}
		}

		<-done
		stopServer(serverProcess)
	}
}
