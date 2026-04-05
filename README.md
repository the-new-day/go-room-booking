# Room Booking Service

REST API сервис для бронирования переговорок: комнаты, расписания, слоты и брони.
Спецификация API находится в `api.yaml`.

Реализованные дополнительные задания:
* Регистрация и авторизация по email/паролю
* Makefile

## Стек

- Go 1.25
- PostgreSQL 16
- chi (роутинг)
- pgx/v5 (драйвер БД)
- golang-jwt (JWT)
- go-playground/validator (валидация)

## Быстрый старт

```bash
docker compose up --build
```

Сервис поднимается на `http://localhost:8080`.

Миграции:

```bash
make migrate-up
```

Сид данных (dummy users):

```bash
make seed
```

## Команды Makefile

| Команда | Назначение |
|---|---|
| `make up` | Запуск сервиса и зависимостей |
| `make down` | Остановка |
| `make migrate-up` | Применить миграции |
| `make migrate-down` | Откатить последнюю миграцию |
| `make migrate-down-all` | Откатить все миграции |
| `make migrate-version` | Показать версию миграций |
| `make seed` | Наполнить БД тестовыми пользователями |
| `make test` | Запустить тесты |

## Авторизация

`POST /dummyLogin` принимает `{"role":"admin"}` или `{"role":"user"}` и возвращает JWT.

Фиксированные UUID для тестирования:
- admin: `00000000-0000-0000-0000-000000000001`
- user: `00000000-0000-0000-0000-000000000002`

Токен передаётся в заголовке:

```
Authorization: Bearer <token>
```

## Генерация слотов

Слоты создаются при создании расписания. Сейчас генерируется окно на **сегодня + 7 дней вперёд** (всего 8 суток).
Длительность слота фиксирована: 30 минут. Даты и время — только UTC.

## Тесты

### Юнит‑тесты
```bash
make test
```

### E2E
Тесты находятся в `tests/e2e`. Они обращаются к запущенному сервису.

```bash
go test ./tests/e2e -v
```

Если сервис доступен на другом адресе:

```bash
E2E_BASE_URL=http://localhost:8080 go test ./tests/e2e -v
```

## Переменные окружения

| Переменная | По умолчанию | Описание |
|---|---|---|
| `HTTP_PORT` | `8080` | Порт сервиса |
| `POSTGRES_HOST` | `localhost` | Хост БД |
| `POSTGRES_PORT` | `5432` | Порт БД |
| `POSTGRES_USER` | `user` | Пользователь БД |
| `POSTGRES_PASSWORD` | `pass` | Пароль БД |
| `POSTGRES_DB` | `meeting_booking` | База данных |
| `POSTGRES_SSLMODE` | `disable` | SSL‑режим |
| `JWT_SIGN_KEY` | `secret_key` | Ключ подписи JWT |
