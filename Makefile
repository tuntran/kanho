.PHONY: build generate migrate test lint dev clean

build:
	cd web && npm run build
	go build -o bin/kanho ./cmd/server

generate:
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen \
		--config api/oapi-codegen.yaml api/openapi.yaml
	cd web && npx @hey-api/openapi-ts \
		--input ../api/openapi.yaml \
		--output src/api/generated

migrate:
	go run ./cmd/server

test:
	go test ./...

lint:
	golangci-lint run

dev:
	docker compose up --build

clean:
	rm -rf bin/ web/dist/
