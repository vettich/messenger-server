run: gen
	go run server.go

build: gen
	go build -o bin/server server.go

gen:
	go generate
