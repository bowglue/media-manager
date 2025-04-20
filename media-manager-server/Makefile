# Docker commands
build-server-base:
	@docker build -f Dockerfile.server.base -t media-manager-server-base:latest .

build-service-base:
	@docker build -f ./services/Dockerfile.service.base -t media-manager-service-base:latest .

build:
	@docker-compose build

start:
	@docker-compose up -d

stop:
	@docker-compose stop

down:
	@docker-compose down

bash-gateway:
	@docker exec -it gateway /bin/sh

bash-user-service:
	@docker exec -it user-service /bin/sh

bash-streaming-service:
	@docker exec -it streaming-service /bin/sh

logs:
	@docker-compose logs -f

generate-proto:
	@docker-compose run --rm proto-gen

#protoc --go_out=paths=import:proto/generated --go-grpc_out=paths=import:proto/generated proto/definitions/*.proto


