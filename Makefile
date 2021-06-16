ifneq (,$(wildcard ./.env))
    include .env
    export
endif

build: fmt lint test
	go build -o ./bin/storeservice cmd/main.go
	docker-compose -f docker/docker-compose.yml build

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

migrate_up:
	migrate -database "$(STORE_DATABASE_DRIVER)://$(STORE_DATABASE_USER):$(STORE_DATABASE_PASSWORD)@tcp(localhost:3371)/$(STORE_DATABASE_NAME)" -path ./migrations up

migrate_down:
	migrate -database "$(STORE_DATABASE_DRIVER)://$(STORE_DATABASE_USER):$(STORE_DATABASE_PASSWORD)@tcp(localhost:3371)/$(STORE_DATABASE_NAME)" -path ./migrations down -all