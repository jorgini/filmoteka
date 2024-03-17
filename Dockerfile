FROM golang
LABEL authors="jorgini"

WORKDIR /usr/src

COPY["go.mod", "go.sum", "./"]
RUN go mod download

EXPOSE 3030

COPY app ./

CMD["go run cmd/main.go"]