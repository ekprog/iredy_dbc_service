FROM --platform=linux/amd64 ubuntu:20.04

# Creating root dir
WORKDIR /go/src/microservice/

COPY ./ /go/src/microservice/

RUN apt update
RUN apt-get -y upgrade
RUN apt -y install wget
RUN wget https://go.dev/dl/go1.20.2.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf ./go1.20.2.linux-amd64.tar.gz
RUN rm go1.20.2.linux-amd64.tar.gz

# EXPORT
ENV GOROOT=/usr/local/go
ENV PATH=$PATH:/usr/local/go/bin

RUN apt install -y protobuf-compiler
RUN apt-get -y install golang-goprotobuf-dev
RUN go get -u google.golang.org/protobuf

ENV PATH=$PATH:/$GO_PATH/bin

RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
RUN go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
ENV PATH="$PATH:$(go env GOPATH)/bin"
RUN export PATH="$PATH:$(go env GOPATH)/bin"

RUN protoc -I ./proto \
--go_out ./pkg/pb \
--go_opt paths=source_relative \
./proto/events/*.proto

RUN protoc -I ./proto \
--go_out ./pkg/pb \
--go_opt paths=source_relative \
--go-grpc_out ./pkg/pb \
--go-grpc_opt paths=source_relative \
--grpc-gateway_out ./pkg/pb \
--grpc-gateway_opt paths=source_relative \
./proto/api/*.proto

RUN go get .
RUN go build -o /go-app

EXPOSE 8071

CMD [ "/go-app" ]
