## **Homework: week 4**

## **Task Manager REST API**

Простой REST API-сервис на Go с использованием Fiber. Предоставляет базовые возможности для управления задачами. Хранение данных реализовано в памяти (in-memory).

### **Возможности**
- Создание задач через API
- Валидация входных данных
- Логирование через zap
- Хранение данных в оперативной памяти
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
```
2. Установите зависимости и запустите проект:
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
    "task_id": "cfbfdf02-ea2e-45fa-ae88-3fd81ab939bf"
  }
}
```

### **Получение задачи по ID**
GET /v1/task/{id}
```bash
{
    "status": "success",
    "data": {
        "id": "767a8b31-9cd6-4175-9610-396b0b4aa154",
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
        "id": "767a8b31-9cd6-4175-9610-396b0b4aa154",
        "title": "Updated task",
        "description": "All routes done",
        "status": "done",
        "created_at": "2025-05-16T16:46:52.058644+04:00",
        "updated_at": "2025-05-16T16:47:01.66315+04:00"
    }
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
