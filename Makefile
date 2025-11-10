.PHONY: all clean run run-watch run-server server

all: templ es server

run-server:
	./bin/server

run:
	go run main.go

run-watch:
	go run watch.go

templ:
	templ generate -path "templates/"

es:
	./node_modules/.bin/esbuild src/index.js --bundle --outfile=./static/index.js

server:
	go build -o bin/server server/server.go

clean:
	rm -f templates/*_templ.go
	rm -f static/index.js
	rm -f bin/server
