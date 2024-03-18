FROM golang:latest
LABEL authors="jorgini"

WORKDIR /usr/src

RUN go version
ENV GOPATH=/

RUN apt-get update
RUN apt-get -y install postgresql-client

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY ./ ./


# make wait-for-postgres.sh executable
RUN chmod +x wait_for_postgres.sh

EXPOSE 8000

CMD ["go run cmd/main.go"]