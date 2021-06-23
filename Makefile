ifneq (,$(wildcard ./.env))
    include .env
    export
endif

build: modules fmt proto lint test
	go build -o ./bin/storeservice cmd/http/main.go
	go build -o ./bin/storegrpcservice cmd/grpc/main.go

fmt:
	go fmt ./...

test:
	go test ./...

lint:
	golangci-lint run

db:
	mysql -h 127.0.0.1 -u $(STORE_DATABASE_USER) -p$(STORE_DATABASE_PASSWORD) $(STORE_DATABASE_NAME)

migrate_up:
	migrate -database "$(STORE_DATABASE_DRIVER)://$(STORE_DATABASE_USER):$(STORE_DATABASE_PASSWORD)@tcp(localhost:3371)/$(STORE_DATABASE_NAME)" -path ./migrations up

migrate_down:
	migrate -database "$(STORE_DATABASE_DRIVER)://$(STORE_DATABASE_USER):$(STORE_DATABASE_PASSWORD)@tcp(localhost:3371)/$(STORE_DATABASE_NAME)" -path ./migrations down -all

modules:
	go mod tidy

proto:
	cd ./api && protoc -I/usr/local/include -I. \
                 -I$$GOPATH/src \
                 -I. \
                 -I$$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis \
                 --swagger_out=logtostderr=true:. \
                 --grpc-gateway_out=logtostderr=true:. \
                 --go_out=plugins=grpc:. ./storeservice.proto