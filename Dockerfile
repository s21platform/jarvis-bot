FROM golang:1.24-alpine as builder

WORKDIR /usr/src/service
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN go build -o build/main cmd/bot/main.go
RUN go build -o build/birthday_pool cmd/worker/birthday_pool/main.go
RUN go build -o build/service cmd/service/main.go


FROM alpine:latest

WORKDIR /app

COPY --from=builder /usr/src/service/build/main .
COPY --from=builder /usr/src/service/build/birthday_pool .
COPY --from=builder /usr/src/service/build/service .
CMD "/app/main" & "/app/birthday_pool" & "/app/service"