CURRENT_DIR=$(shell pwd)

APP=$(shell basename ${CURRENT_DIR})
APP_CMD_DIR=${CURRENT_DIR}/cmd

TAG=latest
ENV_TAG=latest

build:
	CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o ${CURRENT_DIR}/bin/${APP} ${APP_CMD_DIR}/main.go

gen-proto-module:
	./scripts/gen_proto.sh ${CURRENT_DIR}

rm-proto-omit-empty:
	chmod 744 ./scripts/rm_omit_empty.sh && ./scripts/rm_omit_empty.sh ${CURRENT_DIR}

swag-init:
	swag init -g api/api.go -o api/docs --parseDependency --parseVendor
run:
	go run cmd/main.go

linter:
	golangci-lint run

migration-up:
	migrate -path ./migrations -database 'postgres://admin:password@0.0.0.0:5454/online_banking?sslmode=disable' up

migration-down:
	migrate -path ./migrations -database 'postgres://admin:password@0.0.0.0:5454/online_banking?sslmode=disable' down

mock:
	mockgen -destination storage/mock/mock.go -source=storage/storage.go

docker-build:
	docker compose up --build -d