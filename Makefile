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

cluster:
	kubectl apply -f deployment/sc.yaml
	kubectl apply -f deployment/pv.yaml
	kubectl apply -f deployment/pvc.yaml
	kubectl apply -f deployment/mongo-config.yaml
	kubectl apply -f deployment/mongo-secret.yaml
	kubectl apply -f deployment/httpserver-config.yaml
	kubectl apply -f deployment/httpserver-secret.yaml
	kubectl apply -f deployment/httpserver-tls-secret.yaml
	kubectl apply -f deployment/mongo.yaml
	kubectl apply -f deployment/httpserver.yaml
	kubectl apply -f deployment/httpserver-ingress.yaml

destroy:
	kubectl delete --ignore-not-found=true -f deployment/httpserver-ingress.yaml
	kubectl delete --ignore-not-found=true -f deployment/httpserver.yaml
	kubectl delete --ignore-not-found=true -f deployment/mongo.yaml
	kubectl delete --ignore-not-found=true -f deployment/httpserver-tls-secret.yaml
	kubectl delete --ignore-not-found=true -f deployment/httpserver-secret.yaml
	kubectl delete --ignore-not-found=true -f deployment/httpserver-config.yaml
	kubectl delete --ignore-not-found=true -f deployment/mongo-secret.yaml
	kubectl delete --ignore-not-found=true -f deployment/mongo-config.yaml
	kubectl delete --ignore-not-found=true -f deployment/pvc.yaml
	kubectl delete --ignore-not-found=true -f deployment/pv.yaml
	kubectl apply -f deployment/sc.yaml