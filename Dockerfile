FROM golang:1.19.0

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o urlshort ./cmd/server/main.go

CMD ["./urlshort"]