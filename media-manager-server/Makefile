# Docker commands
build:
	@docker-compose -f docker-compose.yml build

start:
	@docker-compose -f docker-compose.yml up -d

stop:
	@docker-compose -f docker-compose.yml stop

down:
	@docker-compose -f docker-compose.yml down
	
bash:
	@docker exec -it media-manager-server /bin/sh

logs:
	@docker-compose logs -f

# Database file location
DB_FILE ?= ./database/data/mydb.db

# SQL commands
generate-sql:
	@echo "Cleaning up old repository code..."
	@rm -rf internal/repository/
	@echo "Repository code cleanup complete."
	@echo "Generating SQL code..."
	sqlc generate -f database/sqlc.yaml
	@echo "Code generation complete."

# Migration commands with automatic code generation
migrate-up:
	migrate -path ./database/migrations -database "sqlite3://$(DB_FILE)" up
	$(MAKE) generate-sql
	@echo "Migrations applied and code generated."

migrate-down:
	migrate -path ./database/migrations -database "sqlite3://$(DB_FILE)" down 1
# $(MAKE) generate-sql
# @echo "Migration reverted and code regenerated."

# Create new migration files
migrate-create:
	@read -p "Enter migration name: " name; 
	@migrate create -ext sql -dir ./database/migrations -seq $$name
