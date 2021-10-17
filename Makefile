dep:
	go mod vendor

build: dep
	go build -o cncamp_http_server main.go

run: dep
	go run main.go

test: dep
	go test cncamp_a01/api -v
	go test cncamp_a01/config -v
	go test cncamp_a01/constant -v