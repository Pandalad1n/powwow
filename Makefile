IMG = powwow

.PHONY: build
build:
	docker build -t ${IMG}_serever -f ${PWD}/cmd/server/Dockerfile .
	docker build -t ${IMG}_client -f ${PWD}/cmd/client/Dockerfile .
	docker network create ${IMG} | true

.PHONY: start_client
start_client:
	docker run --rm --net=${IMG} --name=${IMG}_client ${IMG}_client

.PHONY: start_server
start_server:
	docker run --rm  --net=${IMG} --name=${IMG}_server ${IMG}_serever

.PHONY: test
test:
	go test ./... -race -timeout 2m

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
