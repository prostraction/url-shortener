FROM golang:1.19.0
RUN apt-get update
RUN apt install -y protobuf-compiler
RUN go version
ENV GOPATH=/
COPY ./ ./
RUN go mod download
RUN go build -o ./bin/server ./cmd/server/main.go
RUN go build -o ./bin/client ./cmd/client/main.go

CMD ["./bin/server"]