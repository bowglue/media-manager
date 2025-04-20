# Docker commands
build-server-base:
	@docker build -f Dockerfile.server.base -t media-manager-server-base:latest .

build-service-base:
	@docker build -f ./services/Dockerfile.service.base -t media-manager-service-base:latest .

build-proto-gen:
	@docker build -t media-manager-proto-generator ./shared/proto 

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

proto-gen:
ifeq ($(FILES),all)
	@echo "Compiling all .proto files..."
	@cd shared/proto/definitions && for %%f in (*.proto) do ( \
		docker run --rm -v "$(CURDIR)/shared/proto:/proto" \
		media-manager-proto-generator \
		--proto_path=/proto \
		--go_out=/proto/generated \
		--go-grpc_out=/proto/generated \
		definitions/%%f )
else
	@echo "Compiling specified proto files: $(FILES)"
	@for %%f in ($(FILES)) do ( \
		docker run --rm -v "$(CURDIR)/shared/proto:/proto" \
		media-manager-proto-generator \
		--proto_path=/proto \
		--go_out=/proto/generated \
		--go-grpc_out=/proto/generated \
		definitions/%%f \
	)
endif


