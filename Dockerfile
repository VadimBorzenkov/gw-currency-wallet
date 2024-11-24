# Этап сборки
FROM golang:1.22-alpine3.19 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы для установки зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Сборка приложения
RUN go build -o main ./cmd/main.go

# Этап выполнения
FROM alpine:3.19

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем бинарный файл из стадии сборки
COPY --from=builder /app/main .

# Копируем необходимые файлы
COPY certs ./certs/

# Открываем порт для работы приложения
EXPOSE 8080

# Команда для запуска приложения
CMD ["./main"]
