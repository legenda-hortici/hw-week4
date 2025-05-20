## **Homework: week 4**

## **Task Manager REST API**

Простой REST API-сервис на Go с использованием Fiber. Предоставляет базовые возможности для управления задачами. Хранение данных реализовано в базе данных (PostgreSQL).

### **Возможности**
- Создание задач через API
- Валидация входных данных
- Логирование через zap
- Хранение данных в БД
- Загрузка конфигурации из .env

### **Настройка проекта**
1. Создайте файл .env в корне проекта и добавьте следующие параметры:
```bash
# Общие настройки приложения
LOG_LEVEL=info

# Настройки REST API
PORT=8080
WRITE_TIMEOUT=15s
SERVER_NAME=SimpleService
TOKEN=123

# Настройки базы данных
DB_HOST=localhost
DB_PORT=5432
DB_USER=ваш_пользователь
DB_PASSWORD=ваш_пароль
DB_NAME=ваша_база_данных
DB_SSL_MODE=disable
DB_POOL_MAX_CONNS=10
DB_POOL_MAX_CONN_LIFETIME=300s
DB_POOL_MAX_CONN_IDLE_TIME=150s
```
2. Создайте контейнер в Docker с базой данных PostgreSQL:
```bash
docker run --name your-container-name -e POSTGRES_USER=ваш_пользователь -e POSTGRES_PASSWORD=ваш_пароль -e POSTGRES_DB=ваша_база_данных -p 5432:5432 -d postgres:latest
```

3. Выполните миграции базы данных:


3. Установите зависимости и запустите проект:
```bash
go run cmd/main.go
```
Сервис будет доступен по адресу: http://localhost:8080

##Примеры запросов
Не забудьте указать заголовок Authorization: Bearer your_secret_token в каждом запросе.

### **Создание задачи**
POST /v1/tasks
```bash
{
  "title": "New Feature",
  "description": "Develop new API endpoint",
  "status": "new"
}
```

### Ответ:
```bash
{
  "status": "success",
  "data": {
    "task_id": "1"
  }
}
```

### **Получение задачи по ID**
GET /v1/task/{id}
```bash
{
    "status": "success",
    "data": {
        "id": "1",
        "title": "Updated task",
        "description": "All routes done",
        "status": "done",
        "created_at": "2025-05-16T16:46:52.058644+04:00",
        "updated_at": "2025-05-16T16:47:01.66315+04:00"
    }
}
```

### **Получение всех задач**
GET /v1/tasks
```bash
{
    "status": "success",
    "data": {
        "id": "1",
        "title": "Updated task",
        "description": "All routes done",
        "status": "done",
        "created_at": "2025-05-16T16:46:52.058644+04:00",
        "updated_at": "2025-05-16T16:47:01.66315+04:00"
    }
}
```

### **Получение всех задач c пагинацией**
```bash
GET /v1/tasks?page=2
{
    "status": "success",
    "data": [
        {
            "id": 3,
            "title": "3 Task in tasklist",
            "description": "Testing RESTAPI service with database",
            "status": "new",
            "created_at": "2025-05-20T14:38:50.207828Z",
            "updated_at": "2025-05-20T14:38:50.207828Z"
        },
        {
            "id": 4,
            "title": "4 Task in tasklist",
            "description": "Testing pagination",
            "status": "new",
            "created_at": "2025-05-20T14:48:11.829368Z",
            "updated_at": "2025-05-20T14:48:11.829368Z"
        }
    ]
}
```

### **Удаление задачи**
DELETE /v1/delete/{id}
```bash
{
  "status": "success"
}
```

### **Обновление задачи**
PUT /v1/update/{id}
```bash
{
  "title": "test1",
  "description": "test1",
  "status": "in_progress"
}
```
### Ответ:
```bash
{
  "status": "success"
}
```

Сервис готов к использованию и может служить основой для полноценного приложения с хранением задач.
