run: gen
	go run server.go

build:
	go build -o bin/server server.go

gen:
	go generate
