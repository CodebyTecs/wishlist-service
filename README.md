# Wishlist Service

REST API сервис для вишлистов: пользователь регистрируется, создает вишлисты, добавляет подарки, делится публичной ссылкой и позволяет резервировать подарки без авторизации.

## Стек

- Go 1.24
- PostgreSQL 16
- Docker Compose
- JWT (HS256)
- Goose (миграции)
- Swaggo (OpenAPI/Swagger)

## Что реализовано

### Обязательная часть

- `POST /auth/register`
- `POST /auth/login`
- `GET /users/me`
- `POST /wishlists`
- `GET /wishlists`
- `GET /wishlists/{id}`
- `PATCH /wishlists/{id}`
- `DELETE /wishlists/{id}`
- `POST /wishlists/{wishlistID}/items`
- `GET /wishlists/{wishlistID}/items`
- `GET /wishlists/{wishlistID}/items/{itemID}`
- `PATCH /wishlists/{wishlistID}/items/{itemID}`
- `DELETE /wishlists/{wishlistID}/items/{itemID}`
- `GET /public/{token}`
- `POST /public/{token}/reserve/{itemID}`

### Дополнительно

- `GET /healthz`
- Swagger аннотации и генерация OpenAPI (`make swagger`)
- Swagger UI: `GET /swagger/index.html`
- Unit-тесты для бизнес-логики (`internal/service`)
- Graceful shutdown
- Автоприменение миграций при старте приложения

## Архитектура

Проект организован в стиле clean architecture:

- `internal/domain` — доменные сущности, входные модели use-case, доменные ошибки
- `internal/service` — бизнес-логика
- `internal/repository` — PostgreSQL реализации интерфейсов репозиториев
- `internal/handlers` — HTTP handlers
- `internal/adapters/http` — router, middleware, DTO, единый формат ошибок
- `internal/adapters/jwt` — генерация и парсинг JWT
- `internal/app` — bootstrap/wiring зависимостей, запуск сервера
- `internal/config` — конфигурация из env
- `migrations` — SQL-миграции Goose
- `docs/swagger` — сгенерированная OpenAPI документация
- `cmd/wishlist-service` — точка входа приложения

## Ключевые решения

1. Все приватные операции по вишлистам и подаркам идут с owner-check.
2. Публичный доступ к вишлисту реализован по `public_token`.
3. Резервирование подарка сделано race-safe через атомарный `UPDATE ... WHERE is_reserved = false`.
4. Ошибки API возвращаются в едином JSON формате:

```json
{
  "error": "invalid request"
}
```

## Запуск

### Одной командой

```bash
docker-compose up --build
```

Сервис будет доступен на `http://localhost:8080`.

Остановить:

```bash
docker-compose down -v
```

Также можно через Makefile:

```bash
make up
make down
```

## Конфигурация

Используются переменные окружения. Пример:

- `.env.example`

Основные переменные:

- `HTTP_SERVER_ADDRESS`
- `HTTP_SERVER_PORT`
- `HTTP_SERVER_TIMEOUT`
- `DB_HOST`
- `DB_PORT`
- `DB_NAME`
- `DB_USER`
- `DB_PASSWORD`
- `DB_SSLMODE`
- `JWT_SECRET`
- `JWT_TTL`

## Миграции

При старте приложения миграции применяются автоматически.

Ручные команды (опционально):

```bash
make migrate-status
make migrate-up
make migrate-down
```

## Swagger

Сгенерировать OpenAPI:

```bash
make swagger
```

Результат:

- `docs/swagger/docs.go`
- `docs/swagger/swagger.json`
- `docs/swagger/swagger.yaml`

Swagger UI:

- `http://localhost:8080/swagger/index.html`

## Примеры запросов (curl)

```bash
BASE_URL=http://localhost:8080
EMAIL=user@example.com
PASSWORD=qwerty123
```

### 1) Регистрация

```bash
curl -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "'"$EMAIL"'",
    "password": "'"$PASSWORD"'"
  }'
```

### 2) Логин

```bash
curl -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "'"$EMAIL"'",
    "password": "'"$PASSWORD"'"
  }'
```

Ответ содержит:

```json
{
  "access_token": "<jwt>"
}
```

```bash
TOKEN="<jwt>"
```

### 3) Текущий пользователь

```bash
curl "$BASE_URL/users/me" \
  -H "Authorization: Bearer $TOKEN"
```

### 4) Создать вишлист

```bash
curl -X POST "$BASE_URL/wishlists" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Birthday 2026",
    "description": "My gifts",
    "event_date": "2026-12-10"
  }'
```

Из ответа нужно взять:

- `id` (wishlistID)
- `public_token`

### 5) Добавить подарок в вишлист

```bash
WISHLIST_ID="<wishlist_id>"

curl -X POST "$BASE_URL/wishlists/$WISHLIST_ID/items" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Mechanical Keyboard",
    "description": "75% layout",
    "product_link": "https://example.com/item",
    "priority": 5
  }'
```

### 6) Публичный просмотр вишлиста

```bash
PUBLIC_TOKEN="<public_token>"

curl "$BASE_URL/public/$PUBLIC_TOKEN"
```

### 7) Публичное резервирование подарка

```bash
ITEM_ID="<item_id>"

curl -X POST "$BASE_URL/public/$PUBLIC_TOKEN/reserve/$ITEM_ID"
```

Повторный резерв того же `item` вернет `409`:

```json
{
  "error": "item already reserved"
}
```

## Проверки качества

Тесты:

```bash
make test
```

Покрытие:

```bash
make test-cover
```

Линтер:

```bash
make lint
```
