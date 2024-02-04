.Phony: build
build: prep
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/site/netlify/functions/. ./cmds/parseopml
	cp -r web/* build/site/.

.Phony: clean
clean:
	-rm -r build

.Phony: prep
prep: clean
	mkdir -p build/site
	mkdir -p build/site/netlify/functions

.Phony: run
run:
	( cd web ; python3 -m http.server )