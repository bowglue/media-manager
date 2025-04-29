# Docker commands
build-server-base:
	@docker build -f Dockerfile.server.base -t media-manager-server-base:latest .

build-service-base:
	@docker build -f ./services/Dockerfile.service.base -t media-manager-service-base:latest .

build-proto-gen:
	@docker build -t media-manager-proto-generator ./shared/proto 

build:
	@docker-compose -f docker-compose.dev.yml build

start:
	@docker-compose -f docker-compose.dev.yml up -d

stop:
	@docker-compose -f docker-compose.dev.yml stop

down:
	@docker-compose -f docker-compose.dev.yml down
	
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
	@cd shared/proto/protobuf && for %%f in (*.proto) do ( \
		docker run --rm -v "$(CURDIR)/shared/proto:/proto" \
		media-manager-proto-generator \
		--proto_path=/proto/protobuf \
		--go_out=/proto/api \
		--go-grpc_out=/proto/api \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		%%f )
else
	@echo "Compiling specified proto files: $(FILES)"
	@for %%f in ($(FILES)) do ( \
		docker run --rm -v "$(CURDIR)/shared/proto:/proto" \
		media-manager-proto-generator \
		--proto_path=/proto/protobuf \
		--go_out=/proto/api \
		--go-grpc_out=/proto/api \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		%%f \
	)
endif


