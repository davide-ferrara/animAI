package main

import (
	"log"
	"os"
	"os/exec"
	"strings" // Necessario per il filtraggio

	// ... altre importazioni
	"github.com/fsnotify/fsnotify"
)

func Watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Canale per il blocco (già presente)
	done := make(chan struct{})

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// 1. FILTRO: Ignora i file generati
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
						// B. MODIFICA CODICE SERVER (GO/TEMPL): Compilazione Templ + Riavvio Server
						log.Printf("Server code change detected in: %s. Recompiling Templ...", event.Name)

						// 1. Compilazione TEMPL
						compileTempl := exec.Command("make", "templ")
						compileTempl.Stdout = os.Stdout
						compileTempl.Stderr = os.Stderr

						if err := compileTempl.Run(); err != nil {
							log.Println("ERROR executing make templ:", err)
							continue // Non riavviare se la compilazione fallisce
						}
						log.Println("Templ Recompiled!")

						// 2. RIAVVIO DEL SERVER (kill the main thread)
						log.Println("Restarting Go application...")
						// Invece di riavviare, qui dovresti inviare un segnale al
						// processo principale che si occupa del server HTTP.
						// Per un setup rapido, puoi usare un tool esterno come 'air' o 'gow'
						// che è progettato per gestire il riavvio del processo genitore.

						// Per il tuo attuale setup: devi uscire dalla watch-loop e innescare
						// il riavvio (il che è complesso senza un gestore esterno).
						// La soluzione più semplice è usare un tool di live-reload esterno,
						// o chiudere il canale per far terminare il programma Go.

						// Per simulare il riavvio: chiudi il watcher e il canale 'done'
						// Questo in un sistema reale verrebbe gestito da un tool di terze parti.
						watcher.Close()
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

	// Aggiungi la directory e blocca
	if err = watcher.Add("."); err != nil {
		log.Fatal(err)
	}
	// Blocca finché la goroutine non chiude 'done' (solo in caso di riavvio)
	<-done
	// Quando 'done' è chiuso, l'esecuzione esce da Watch() e main() termina.
}
