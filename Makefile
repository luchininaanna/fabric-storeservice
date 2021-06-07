ifneq (,$(wildcard ./docker/.env))
    include docker/.env
    export
endif

build: fmt lint test
	go build -o ./bin/storeservice cmd/storeservice/main.go

fmt:
	go fmt ./...

test:
	go test ./...

lint:
	golangci-lint run

up:
	docker-compose -f docker/docker-compose.yml up -d

down:
	docker-compose -f docker/docker-compose.yml down

db:
	mysql -h 127.0.0.1 -u $(STORE_DATABASE_USER) -p$(STORE_DATABASE_PASSWORD) $(STORE_DATABASE_NAME)