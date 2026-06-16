# Web Service

Легкий RESTful веб-сервис, построенный на Go с использованием фреймворка Gin и базы данных SQLite.

## Структура проекта

```
web-service/
├── cmd/
│   └── main.go                 # Точка входа приложения
├── config/
│   └── local.yaml              # Файл локальной конфигурации
├── internal/
│   ├── config/
│   │   └── config.go           # Загрузка и парсинг конфигурации
│   ├── httpserver/
│   │   └── httpserver.go       # Настройка HTTP-сервера и маршруты
│   ├── memory/
│   │   └── memory.go           # Реализация хранилища в памяти
│   ├── models/
│   │   └── models.go           # Модели данных и структуры
│   └── storage/
│       └── sqlite/
│           └── sqlite.go       # Реализация базы данных SQLite
├── storage/                    # Директория для файлов базы данных
├── .gitignore
├── go.mod                      # Зависимости Go модуля
└── README.md                   # Этот файл
```

## Технологический стек

- **Язык программирования**: Go 1.25.0
- **Web-фреймворк**: Gin (v1.12.0)
- **База данных**: SQLite (mattn/go-sqlite3)
- **Конфигурация**: cleanenv (ilyakaznacheev/cleanenv)

## Требования

- Go 1.25.0 или выше
- SQLite3

## Установка

Клонируйте репозиторий и установите зависимости:

```bash
go mod download
go mod tidy
```

## Конфигурация

Конфигурация загружается из файла `config/local.yaml`. Обновите этот файл со своими параметрами перед запуском приложения.

## Запуск приложения

### Вариант 1: Запуск с помощью Go

```bash
go run ./cmd/main.go
```

### Вариант 2: Сборка и запуск бинарного файла

Соберите исполняемый файл:

```bash
go build -o web-service.exe ./cmd/main.go
```

Запустите бинарный файл:

```bash
./web-service.exe
```

Или на Linux/Mac:

```bash
./web-service
```

## Компоненты проекта

### cmd/main.go
Точка входа приложения, инициализирует конфигурацию, базу данных и HTTP-сервер.

### internal/config/config.go
Обрабатывает загрузку конфигурации из YAML-файлов с использованием cleanenv.

### internal/httpserver/httpserver.go
Настраивает HTTP-сервер Gin и определяет все API маршруты и обработчики.

### internal/models/models.go
Определяет структуры данных, используемые во всем приложении.

### internal/storage/
Реализует слой хранилища данных с двумя вариантами:
- **memory.go**: Хранилище в памяти
- **sqlite/sqlite.go**: Персистентное хранилище SQLite

## Разработка

Директория `internal/` содержит весь код приложения, организованный по функциональности, следуя лучшим практикам Go для структуры проекта.

## Функциональность

1. Создание короткой ссылки

POST /links

Request

{
  "url": "https://example.com/some/very/long/url"
}

Response

{
  "short_code": "abc123"
}


2. Получение оригинальной ссылки

GET /links/{short_code}

Response

{
  "url": "https://example.com/some/very/long/url",
  "visits": 15
}



3. Получение списка ссылок

GET /links

Поддержка пагинации.

Query параметры:

limit  
offset


4. Удаление ссылки

DELETE /links/{short_code}


5. Статистика

GET /links/{short_code}/stats

Response

{
  "short_code": "abc123",
  "url": "https://example.com",
  "visits": 15,
  "created_at": "timestamp"
}

