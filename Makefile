dep:
	cd httpserver && go mod vendor

run: dep
	cd httpserver && go run main.go

test: dep
	cd httpserver && go test ./... -v

bin: dep
	cd httpserver && go build -o ../bin/cncamp_http_server

image:
	docker build -t niennienzz/cncamp_http_server:latest .

push: image
	docker push niennienzz/cncamp_http_server:latest
