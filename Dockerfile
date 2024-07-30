FROM golang:1.22.5-alpine3.20 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o server ./cmd/server


FROM alpine:latest as runner

WORKDIR /root/

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]