FROM golang:1.23.1-alpine AS builder
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download -x

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/server ./cmd/server/main.go

FROM alpine:latest AS final
WORKDIR /root/

COPY --from=builder /bin/server /bin/

EXPOSE 8080
ENTRYPOINT [ "/bin/server" ]