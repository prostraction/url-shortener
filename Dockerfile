FROM golang:1.19.0
RUN go version
ENV GOPATH=/
COPY ./ ./

# install psql
RUN apt-get update
RUN apt install -y postgresql-client

# wait wait-for-postgres.sh executable
RUN chmod +x scripts/wait-for-postgres.sh

# build
RUN go mod download
RUN go build -o ./bin/server ./cmd/server/main.go
RUN go build -o ./bin/client ./cmd/client/main.go

ENTRYPOINT [ "./bin/server" ]