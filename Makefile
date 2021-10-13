all: gen

gen:
	go generate ./pkg/...

build:
	go build ./pkg/...