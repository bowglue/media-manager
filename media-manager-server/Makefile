# Docker commands
build:
	@docker-compose build

start:
	@docker-compose up -d

stop:
	@docker-compose stop

down:
	@docker-compose down
	
bash:
	@docker exec -it media-manager-server /bin/sh

logs:
	@docker-compose logs -f

graph:
	@gqlgen generate

# Database file location
DB_FILE ?= ./data/mm.db

sql:
	@sqlc generate -f sqlc.yaml

migrate-up:
	@migrate -path ./common/database/migrations -database "sqlite3://$(DB_FILE)" up

migrate-down:
	@migrate -path ./common/database/migrations -database "sqlite3://$(DB_FILE)" down 1

