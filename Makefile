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
