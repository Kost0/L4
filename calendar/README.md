# Calendar API

Простое REST API приложение для управления календарными событиями с системой уведомлений.

## Возможности

- Создание, обновление и удаление событий
- Получение событий за день, неделю или месяц
- Автоматические уведомления за 2 часа до начала события
- Хранение данных в PostgreSQL

## Требования

- Docker
- Docker Compose

## Запуск

1. Клонируйте репозиторий:
```bash
git clone https://github.com/Kost0/L4.git
cd L4/calendar
```

2. Запустите приложение с помощью Docker Compose:
```bash
docker-compose up -d
```

Сервер будет доступен по адресу `http://localhost:8080`

## API Endpoints

### Создание события
```http
POST http://localhost:8080/create_event
Content-Type: application/json

{
  "user_id": "ecec1151-c972-433f-9861-6200440ede46",
  "event": "Встреча с клиентом",
  "date": "2025-11-18T14:00:00Z"
}
```

### Обновление события
```http
POST http://localhost:8080/update_event?id={event_id}
Content-Type: application/json

{
  "user_id": "ecec1151-c972-433f-9861-6200440ede46",
  "event": "Обновленное название",
  "date": "2025-11-18T15:00:00Z"
}
```

### Получение событий за день
```http
GET http://localhost:8080/events_for_day?user_id=ecec1151-c972-433f-9861-6200440ede46&date=2025-11-18T08:00:00Z
```

### Получение событий за неделю
```http
GET http://localhost:8080/events_for_week?user_id=ecec1151-c972-433f-9861-6200440ede46&date=2025-11-18T08:00:00Z
```

### Получение событий за месяц
```http
GET http://localhost:8080/events_for_month?user_id=ecec1151-c972-433f-9861-6200440ede46&date=2025-11-18T08:00:00Z
```

### Удаление события
```http
POST http://localhost:8080/delete_event?id={event_id}
```

## Примеры использования

Дополнительные примеры запросов можно найти в файле `http/examples.http`.

## Технологии

- Go
- PostgreSQL
- Docker