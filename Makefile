all: gen

gen:
	go run cmd/gen/*
	mv cmd/gen/token_meta.json pkg/token/token_meta.json

build:
	go build ./pkg/...