.PHONY: all clean run run-watch run-server server

build: templ es server

server:
	go build -o webserver server/server.go

templ:
	templ generate -path "templates/"

es:
	./node_modules/.bin/esbuild src/index.js --bundle --outfile=./static/bundle.js

run:
	go run .

clean:
	rm -f templates/*_templ.go
	rm -f static/index.js
	rm -f bin/server
