FROM golang:1.20-alpine

WORKDIR /app

RUN go install github.com/cosmtrek/air@v1.42.0
RUN go install github.com/roerohan/wait-for-it@v0.2.13

COPY go.mod go.sum ./
RUN go mod download

CMD ["air", "-c", ".air.toml", "serve"]
