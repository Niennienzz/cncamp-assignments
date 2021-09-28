dep:
	go mod vendor

build: dep
	go build -o http_server main.go

run: dep
	go run main.go
