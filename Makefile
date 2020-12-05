GOTAGS :=
MAKEFLAGS += -j2

.PHONY: build bin/white-elephant dev bin/index.html

# .NOTPARALLEL: build
build: bin/dist bin/white-elephant

dev: override GOTAGS += dev
dev: _web-watch _server

_web-watch:
	yarn run parcel watch ./web/index.html --public-url ./static/

_server: bin/white-elephant
	./bin/white-elephant serve

bin/dist:
	yarn run parcel build ./web/index.html --public-url ./static/

bin/white-elephant:
	go generate ./...
	go build --tags='$(GOTAGS)' -o $@ main.go

