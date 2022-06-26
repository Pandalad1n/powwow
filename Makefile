IMG = powwow
IMG_CLIENT = powwow-client
IMG_SERVER = powwow-server
NETWORK = powwow

.PHONY: build
build:
	docker build -t ${IMG_SERVER} -f ${PWD}/cmd/server/Dockerfile .
	docker build -t ${IMG_CLIENT} -f ${PWD}/cmd/client/Dockerfile .
	docker network create ${NETWORK} | true

.PHONY: start-client
start-client:
	docker run --rm --net=${NETWORK} --name=${IMG_CLIENT} ${IMG_CLIENT} /app/client -host ${IMG_SERVER} -port 8888

.PHONY: start-server
start-server:
	docker run --rm --net=${NETWORK} --name=${IMG_SERVER} ${IMG_SERVER} /app/server -port 8888

.PHONY: test
test:
	docker run -it --rm -w /app -v ${PWD}:/app golang:1.17 go test ./... -race -timeout 2m

.PHONY: lint
lint:
	docker run \
		--rm \
		-ti \
		-w /app \
		-v $(PWD):/app \
		golangci/golangci-lint:v1.45-alpine golangci-lint run

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
