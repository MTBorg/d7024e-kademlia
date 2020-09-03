FROM golang:1.15.0-alpine3.12

WORKDIR /kademlia

ENV PATH /kademlia/bin:$PATH

COPY go.mod go.sum ./

# Hot reloading module
RUN go get github.com/githubnemo/CompileDaemon

COPY . .

ENTRYPOINT CompileDaemon \
	--color=true \
	--build="sh scripts/build-dev.sh" \
	--log-prefix=false \
	--command=./bin/kademlia
