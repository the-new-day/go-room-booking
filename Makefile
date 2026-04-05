test:
	go test -v -race ./...
.PHONY: test

up:
	docker compose up --build
.PHONY: up

down:
	docker compose down
.PHONY: down

seed:
	echo "TODO"
.PHONY: seed

migrate-up:
	docker compose --profile migrations run --rm migrator
.PHONY: migrate-up

migrate-down:
	docker compose --profile migrations run --rm migrator down 1
.PHONY: migrate-down

migrate-down-all:
	docker compose --profile migrations run --rm migrator down -all
.PHONY: migrate-down-all

migrate-version:
	docker compose --profile migrations run --rm migrator version
.PHONY: migrate-version

up-with-migrations:
	docker compose --profile migrations run --rm migrator
	docker compose up --build
.PHONY: up-with-migrations
