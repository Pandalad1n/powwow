IMG = powwow

.PHONY: build
build:
	docker build -t ${IMG} .

.PHONY: start_client
start_client:
	docker run --net=host ${IMG} go run ./cmd/client/.

.PHONY: start_server
start_server:
	docker run --net=host ${IMG} go run ./cmd/server/.

.PHONY: gen
gen:
	docker build -t ${IMG} --build-arg ARC=`uname -p` -f proto/Dockerfile .
	docker run \
		--rm \
		-w /app \
		-u `id -u ${USER}` \
		-v ${PWD}/proto:/app \
		${IMG} protoc server.proto \
		   --proto_path=. \
		   --go_opt=paths=source_relative \
		   --go_opt=Mserver.proto=pow/proto \
		   --go_out=.