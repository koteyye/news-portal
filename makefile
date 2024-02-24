USER_PORT := 8080
USER_DSN := postgres://postgres:postgres@localhost:5432/user?sslmode=disable

.PHONY: up
up:
	@docker-compose -f ./scripts/docker-compose.yaml up -d

.PHONY: down
down:
	@docker-compose -f ./scripts/docker-compose.yaml down