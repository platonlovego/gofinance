# GoFinance

REST API для учёта личных финансов — доходы, расходы, статистика по категориям.

## Возможности

- Регистрация и авторизация (JWT)
- CRUD операций с транзакциями (доходы и расходы)
- Категории: еда, транспорт, зарплата и т.д.
- Сводка за период: общий доход, расходы, баланс, разбивка по категориям
- Запуск через Docker Compose

## Стек

Go 1.22 · Gin · PostgreSQL · sqlx · JWT · Docker

## Быстрый старт

```bash
git clone https://github.com/platonlovego/gofinance
cd gofinance
cp .env.example .env
docker-compose up --build
```

API будет доступен на `http://localhost:8080`.

## API

### Авторизация

```
POST /api/v1/auth/register   { "email": "...", "password": "..." }
POST /api/v1/auth/login      { "email": "...", "password": "..." }
```

### Транзакции (требуют заголовок `Authorization: Bearer <token>`)

```
GET    /api/v1/transactions          — список всех транзакций
POST   /api/v1/transactions          — создать транзакцию
GET    /api/v1/transactions/:id      — одна транзакция
PUT    /api/v1/transactions/:id      — обновить
DELETE /api/v1/transactions/:id      — удалить
```

Пример тела запроса:
```json
{
  "category_id": 5,
  "amount": 1200.50,
  "comment": "Пятёрочка",
  "date": "2025-06-15T00:00:00Z"
}
```

### Статистика

```
GET /api/v1/summary?from=2025-06-01&to=2025-06-30
```

```json
{
  "total_income": 85000,
  "total_expense": 34200.50,
  "balance": 50799.50,
  "by_category": [
    { "category_name": "Зарплата", "type": "income", "total": 85000 },
    { "category_name": "Продукты", "type": "expense", "total": 12400 }
  ]
}
```

### Категории

```
GET /api/v1/categories
```

## Структура проекта

```
.
├── cmd/server/main.go          — точка входа
├── config/                     — подключение к БД
├── internal/
│   ├── handler/                — HTTP-обработчики, роутинг, middleware
│   ├── model/                  — структуры данных
│   ├── repository/             — работа с БД
│   └── service/                — бизнес-логика
├── migrations/001_init.sql     — схема БД и стартовые категории
├── Dockerfile
└── docker-compose.yml
```

## Локальный запуск без Docker

```bash
# Запустить только PostgreSQL
docker-compose up db -d

# Применить миграцию вручную
psql -h localhost -U postgres -d gofinance -f migrations/001_init.sql

# Запустить приложение
cp .env.example .env  # указать DB_HOST=localhost
go run ./cmd/server
```

## Планы

- [ ] Экспорт в CSV
- [ ] Повторяющиеся транзакции (подписки)
- [ ] Лимиты по категориям с уведомлением
