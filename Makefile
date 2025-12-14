APP=app

.PHONY: build run test docker-build docker-run tidy

build:
	go build -o $(APP) ./cmd/app

run:
	DB_PATH=./data/devices.db SERVER_ADDR=:8080 ./$(APP)

test:
	go test ./...

tidy:
	go mod tidy

docker-build:
    docker build -t go-backend:latest .

docker-run:
    docker run --rm -p 8080:8080 -v $$(pwd)/data:/data go-backend:latest

.PHONY: postman
postman:
	docker run --rm -v $$(pwd)/docs/postman:/etc/newman postman/newman run /etc/newman/DevicesAPI.postman_collection.json -e /etc/newman/DevicesAPI.postman_environment.json --env-var baseUrl=http://host.docker.internal:8080
