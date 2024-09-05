# Envoy GRPC Trailer Example

Demo of how to grab headers and trailers as part of a gRPC request using the [Envoy](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/ext_proc_filter) proxy. This example is deployed as part of a kubernetes deployment and thus involves building a few yaml files.

The point of this project is to give an idea on how to develop & deploy a very simple gRPC client and server fronted by Envoy with an external processor allowing the interception of the request and response body & headers.

The reason specifically targeting trailers in this project is due to the fact that trailers are usually not thought about, nor even known about most of the time. 

## Deployment

This project contains the following folders:

- client
  - holds `main.go` grpc client and the associated `Dockerfile` to host the binary on an Ubuntu image. A much smaller image can certaily be used due to Go's static linking not requiring much behind glibc.
  - `main.go` should be edited properly to ensure it is correctly targeting the service fronting Envoy/gRPC-server as needed
- server
  - holds `main.go` grpc server and the associated `Dockerfile` built the same way as the server
- externalprocessor
  - holds `main.go` externalprocessor that receives bidirectional gRPC stream messages from Envoy. This processsor stream receives messages depending upon the `envoy.yaml` configuration. It is currently set to be sent request/response headers and trailers, therefore the stream should in total receive 4 messages. It's important to note that the following are from the perspective of the downstream and upstream service, the envoy proxy is transparent in this case.
    - request headers
    - request trailers
    - response headers
    - response trailers
- pb
  - holds associated protobuf Go code
- proto 
  - holds proto definitions
- envoy
  - holds `envoy.yaml` configuration
- yaml
  - holds `kubernetes` yaml files to alter/deploy
- scripts
  - holds useful scripts

**Steps to deploy**:
1. Build the client, server and external processor binaries as below via `make` targets
2. Upload the images to your registry, in testing it was ECR that was used
3. Edit the yaml files to correctly target your uploaded images
4. Deploy the yaml files to your cluster
5. Open up 3 terminals, one for the each client/server/externalprocessor
6. Start the server and externalprocessor binaries
7. Start the client

## Make

`make gen` - generate protobuf files and store in `pb/`

`make clean` - clean `pb/` folders

`make build-client` - build client

`make build-server` - build server

`make build-external` - build external processor

## Success Critiera

If run properly, the externalprocessor should print out correctly request headers, request trailers (usually none), response headers and response trailers to stdout of that process/container. 