FROM golang:latest AS builder

WORKDIR /build

ADD go.mod .

ADD go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bin cmd/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /build .

ENTRYPOINT [ "/app/bin" ]
