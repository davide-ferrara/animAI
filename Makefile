all: templ es run
	echo "Building and running!"

run:
	go run *.go

templ:
	templ generate -path "templates/"
	mv templates/*.go .
es:
	./node_modules/.bin/esbuild src/index.js --bundle --outfile=./static/index.js
