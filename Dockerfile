# Build stage
FROM golang:1.25-alpine AS builder

# Устанавливаем зависимости для компиляции и protoc
RUN apk add --no-cache git protoc protobuf-dev

# Устанавливаем gRPC плагины и добавляем в PATH
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

# Добавляем GOPATH/bin в PATH для protoc плагинов
ENV PATH $PATH:/go/bin

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cometsService/cmd/server

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
RUN adduser -D -s /bin/sh appuser
USER appuser
EXPOSE 8082
CMD ["./main"]