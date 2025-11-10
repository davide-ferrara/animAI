ESBUILD := ./node_modules/.bin/esbuild 
SERVER_BIN := webserver
JS_BUNDLE := static/bundle.js
TEMPLATES := templates/
TEMPLATE_SUFFIX := *_templ.go 

.PHONY: all clean run run-server server

all: build

build: templ es server

server:
	go build -o $(SERVER_BIN) server/server.go

templ:
	templ generate -path "$(TEMPLATES)"

es:
	$(ESBUILD) src/index.js --bundle --outfile=$(JS_BUNDLE)

run:
	go run .

run-server:
	./$(SERVER_BIN)

clean:
	rm -f $(TEMPLATES)$(TEMPLATE_SUFFIX)
	rm -f $(JS_BUNDLE)
	rm -f $(SERVER_BIN)
