.PHONY: docker-up docker-down test

docker-up:
	docker compose up -d --build

docker-down:
	docker compose down

test:
	go test ./...
