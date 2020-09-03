FROM golang:1.15.0-alpine3.12

WORKDIR /kademlia

COPY . .

RUN go build cmd/kademlia/kademlia.go && \
	go build cmd/kademliactl/kademliactl.go

FROM alpine:3.12.0

WORKDIR /kademlia

# Copy binaries
COPY --from=0 /kademlia/kademlia /kademlia/kademliactl ./

# Add the binaries to the path
ENV PATH /kademlia:$PATH

CMD ["kademlia"]
