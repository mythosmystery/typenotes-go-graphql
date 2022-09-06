FROM golang:1.19-alpine AS builder

WORKDIR /usr/src/server

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o server

FROM alpine:latest AS runner

WORKDIR /usr/src/server

COPY --from=builder /usr/src/server/server .

EXPOSE 3001

CMD ["./server"]