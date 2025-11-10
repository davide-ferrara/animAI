package main

import "log"

func main() {
	log.SetPrefix("[MAIN] ")
	log.Println("Starting watch process...")
	Watch()
}
