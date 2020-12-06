GOTAGS :=
MAKEFLAGS += -j2

.PHONY: build bin/white-elephant dev bin/index.html

# .NOTPARALLEL: build
build: bin/dist bin/white-elephant

dev-server: override GOTAGS += dev
dev-server: bin/white-elephant
	./bin/white-elephant serve

parcel-watch:
	yarn run parcel watch ./web/index.html --public-url /static/

cockroach:
	cockroach start-single-node \
		--advertise-addr localhost \
		--insecure \
		--port 4445 \

bin/dist:
	NODE_ENV=production yarn run parcel build ./web/index.html --public-url /static/

bin/white-elephant:
	go generate ./...
	go build --tags='$(GOTAGS)' -o $@ main.go

