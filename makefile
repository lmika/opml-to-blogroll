.Phony: build
build: prep
	go build -o build/site/netlify/functions/. ./cmds/parseopml
	cp -r web/* build/site/.

.Phony: clean
clean:
	-rm -r build

.Phony: prep
prep: clean
	mkdir -p build/site
	mkdir -p build/site/netlify/functions
