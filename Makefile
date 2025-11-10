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
	mv templates/*_templ.go templates/compiled
es:
	./node_modules/.bin/esbuild src/index.js --bundle --outfile=./static/index.js

server:
	go build -o bin/server server/server.go

clean:
	rm -f templates/compiled/*_templ.go
	rm -f static/index.js
	rm -f bin/server
