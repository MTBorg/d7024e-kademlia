FROM golang:1.15.0-alpine3.12

WORKDIR /kademlia

COPY . /kademlia

ENTRYPOINT ["go"]
