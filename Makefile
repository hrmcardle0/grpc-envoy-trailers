PHONY: gen
gen:
	@echo "Generating proto files"
	#protoc -I=proto --go_out=. --go_opt=module=github.com/hrmcardle0/envoy-grpc-trailer --go-grpc_out=. --go-grpc_opt=module=github.com/hrmcardle0/envoy-grpc-trailer proto/*.proto
	protoc -I=proto --go_out=. --go-grpc_out=. proto/*.proto


PHONY: clean
clean: 
	rm pb/*.go
	rm grpc-client/pb/*.go

# build client for ubuntu docker
PHONY: build-client
build-client:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o client/client client/main.go
	docker buildx build --platform linux/amd64 -t envoy-grpc-trailer-client -f client/Dockerfile --load .

# build server for ubuntu docker
PHONY: build-server
build-server:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server/server server/main.go
	docker buildx build --platform linux/amd64 -t envoy-grpc-trailer-server -f server/Dockerfile --load .

# build external process
PHONY: build-external
build-external:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o externalprocessor/externalprocessor externalprocessor/main.go
	docker buildx build --platform linux/amd64 -t envoy-grpc-trailer-externalprocessor -f externalprocessor/Dockerfile --load .
