# Task Service

Это пример реализации REST API для создания, удаления и получения задач

---

## Зависимости

Проект использует минимальный набор зависимостей:

- [Gin](https://github.com/gin-gonic/gin) — HTTP-роутер
- [UUID](https://github.com/google/uuid) — генерация уникальных идентификаторов
- [Zerolog](https://github.com/rs/zerolog) — используемый логгер
- [Testify](https://github.com/stretchr/testify) — фреймворк для тестирования и моков
- [Mockery](https://github.com/golang/mock) — генератор моков (для тестов)

---

## Начало работы

### Установка

1. **Клонируйте репозиторий**:
   ```bash
   git clone https://github.com/prolgrammer/TaskManager
   cd task-service
   ```

2. **Установите зависимости при необходимости**:
   ```bash
   go mod tidy
   ```

3. **Запустите приложение с помощью Docker**:
   ```bash
   docker-compose -f docker-compose.yml up -d
   ```

   Сервис будет доступен по адресу: `http://localhost:8080`.


### Тестирование

Перейдите по адресу http://localhost:8080/swagger/index.html#/default для открытия документации и проведения тестов.

---


## Структура папок

```
task-service/
├── cmd/
│   ├── app/
│   │   └── app.go
│   └── main.go 
├── config/
│   ├── config.go
│   └── config.yaml
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal/
│   ├── controllers/
│   │   ├── http/
│   │   │   ├── middleware/
│   │   │   │   ├── error_handler.go
│   │   │   │   ├── errors.go
│   │   │   │   └── middleware.go
│   │   │   ├── create_task.go
│   │   │   ├── delete_task.go
│   │   │   ├── get_task.go
│   │   │   ├── get_tasks.go
│   │   │   └── router.go
│   │   ├── requests/
│   │   │   └── create_task.go
│   │   ├── responses/
│   │   │   └── task.go
│   │   └── responses/
│   ├── entities/
│   │   └── task.go
│   ├── repositories/
│   │   ├── errors.go
│   │   └── task.go
│   ├── services/
│   │   └── task_manager/
│   │       ├── task_manager.go
│   │       └── worker_pool.go
│   ├── usecases/
│   │   ├── contracts.go
│   │   ├── create_task.go
│   │   ├── delete_task.go
│   │   ├── errors.go
│   │   ├── get_task.go
│   │   └── get_tasks.go
├── pkg/
│   └── logger/
│       ├── logger.go
│       └── zerolog.go
├── .env
├── .env.example
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── README.md
```

### Описание структуры

- **cmd/**: Точка входа приложения.
- **config/**: Конфигурации приложения.
- **internal/controllers/http/**: HTTP-контроллеры, обрабатывающие входящие запросы.
- **internal/controllers/requests/**: Структуры для входных данных.
- **internal/controllers/responses/**: Структуры для ответов.
- **internal/entities/**: Определения сущностей.
- **internal/repositories/**: Репозитории для операций с данными.
- **internal/services**: Вспомогательные сервисы.
- **internal/usecases/**: Бизнес-логика приложения.
- **pkg**: Вспомогательные пакеты.


---


## API

### Создание задачи
- **Метод**: `POST /task`
- **Тело запроса**:
  ```json
  {
    "text": "text"
  }
  ```
- **Ответ** (200 OK):
  ```json
  {
    "task_id": "uuid",
    "status": "created",
    "created_at": "2025-07-06T18:00:00Z",
    "duration": "0s"
  }
  ```
- **Ошибки**:
    - `400 Bad Request`: Некорректный формат запроса.
    - `500 Internal Server Error`: Внутренняя ошибка сервера.
- **Пример**:
  ```bash
  curl -X POST http://localhost:8080/task -H "Content-Type: application/json" -d '{"text":"process_data"}'
  ```

### Получение задачи
- **Метод**: `GET /task/{task_id}`
- **Ответ** (200 OK):
  ```json
  {
    "task_id": "uuid",
    "status": "running",
    "created_at": "2025-07-06T18:00:00Z",
    "duration": "3m45s"
  }
  ```
- **Ошибки**:
    - `404 Not Found`: Задача с указанным `task_id` не найдена.
    - `500 Internal Server Error`: Внутренняя ошибка сервера.
- **Пример**:
  ```bash
  curl http://localhost:8080/task/<task_id>
  ```

### Получение списка задач
- **Метод**: `GET /tasks`
- **Параметры запроса**:
    - `limit` (int): Количество задач на странице (по умолчанию: 10).
    - `offset` (int): Смещение для пагинации (по умолчанию: 0).
- **Ответ** (200 OK):
  ```json
  [
    {
      "task_id": "uuid",
      "status": "running",
      "created_at": "2025-07-06T18:00:00Z",
      "duration": "3m45s"
    },
  ]
  ```
- **Ошибки**:
    - `400 Bad Request`: Некорректный формат запроса.
    - `500 Internal Server Error`: Внутренняя ошибка сервера.
- **Пример**:
  ```bash
  curl http://localhost:8080/tasks?limit=10&offset=0
  ```

### Удаление задачи
- **Метод**: `DELETE /task/{task_id}`
- **Ответ** (200 OK):
  ```json
  {
    "message": "Delete successful"
  }
  ```
- **Ошибки**:
    - `400 Bad Request`: Некорректный формат запроса.
    - `404 Not Found`: Задача с указанным `task_id` не найдена.
    - `500 Internal Server Error`: Внутренняя ошибка сервера.
- **Пример**:
  ```bash
  curl -X DELETE http://localhost:8080/task/<task_id>
  ```


