# Используем образ golang в качестве базового
FROM golang:1.16 AS builder

# Установка переменных окружения
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Создаем директорию внутри образа для размещения исходного кода
WORKDIR /app

# Копируем go.mod и go.sum для загрузки зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости проекта
RUN go mod download

# Копируем исходный код в рабочую директорию
COPY . .

# Собираем бинарный файл приложения
RUN go build -o app

# Второй этап сборки, создаем образ без необходимости устанавливать Go
FROM scratch

# Копируем бинарный файл из предыдущего этапа
COPY --from=builder /app/app /app/app

# Указываем рабочую директорию
WORKDIR /app

# Запускаем приложение
CMD ["./app"]
