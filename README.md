```markdown
# task-manager

REST API для управления задачами на Go. CRUD, фильтрация по статусу, PostgreSQL, юнит- и интеграционные тесты.

[![CI](https://github.com/ImmortaL-jsdev/task-manager/actions/workflows/ci.yml/badge.svg)](https://github.com/ImmortaL-jsdev/task-manager/actions/workflows/ci.yml)

## Стек

- **Go 1.25** — язык
- **gorilla/mux** — роутинг
- **PostgreSQL 16** + **pgx/v5** — БД и драйвер
- **testcontainers-go** — интеграционные тесты
- **gomock** — моки для юнит-тестов
- **golangci-lint** — линтер
- **GitHub Actions** — CI

## Запуск

### Через make

```bash
make docker-up   # поднять PostgreSQL в Docker (порт 5432)
make run         # запустить сервер на :8080
```

### Вручную

```bash
docker compose up -d
go run cmd/server/main.go
```

Проверка, что сервер живой:

```bash
curl http://localhost:8080/health
```

## Конфигурация

Сервер читает настройки БД из переменных окружения. Все имеют дефолтные значения для локального запуска:

| Переменная | По умолчанию | Описание |
|---|---|---|
| `DB_HOST` | `localhost` | Хост PostgreSQL |
| `DB_PORT` | `5432` | Порт |
| `DB_USER` | `task_user` | Пользователь |
| `DB_PASSWORD` | `task_pass` | Пароль |
| `DB_NAME` | `task_db` | Имя БД |

Для продакшена — переопределить через env.

## API

Базовый URL: `http://localhost:8080`

| Метод | Путь | Описание | Успех | Ошибки |
|---|---|---|---|---|
| `GET` | `/health` | Проверка живости (пинг БД) | `200` | `503` |
| `POST` | `/tasks` | Создать задачу | `201` | `400` |
| `GET` | `/tasks` | Список всех задач (свежие первыми) | `200` | — |
| `GET` | `/tasks?status=todo` | Список задач с фильтром по статусу | `200` | `400` |
| `GET` | `/tasks/{id}` | Задача по ID | `200` | `404` |
| `PUT` | `/tasks/{id}` | Обновить задачу (полная замена) | `200` | `400`, `404` |
| `DELETE` | `/tasks/{id}` | Удалить задачу | `204` | `404` |

**Статусы:** `todo`, `in_progress`, `done`. Пустой статус при создании → `todo`.

**Формат ошибки:**

```json
{"error": "title cannot be empty"}
```

**Пример создания:**

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "write tests", "description": "with testcontainers", "status": "todo"}'
```

**Пример обновления (полная замена полей):**

```bash
curl -X PUT http://localhost:8080/tasks/{id} \
  -H "Content-Type: application/json" \
  -d '{"title": "updated", "description": "new desc", "status": "done"}'
```

**Пример удаления:**

```bash
curl -X DELETE http://localhost:8080/tasks/{id}
# → 204 No Content
```

## Тестирование

```bash
make test-unit          # юнит-тесты сервиса (gomock, без Docker)
make test-integration   # интеграционные тесты репозитория (testcontainers + PostgreSQL)
make test-cover         # покрытие всех пакетов
make pre-push           # всё сразу перед пушем: lint + unit + integration
```

Интеграционные тесты требуют **Docker** (testcontainers поднимает `postgres:16-alpine` в контейнере). Переменная `TESTCONTAINERS_RYUK_DISABLED=true` уже прописана в `Makefile` — нужна для работы на Fedora/rootless Docker.

## Структура

```
cmd/server/          — точка входа, конфигурация, graceful shutdown
internal/
  errors/            — кастомные типы ошибок (ValidationError, NotFoundError)
  handlers/          — HTTP-хендлеры
  middleware/        — recovery, logging, CORS
  models/            — модели данных
  repository/        — работа с PostgreSQL (pgx)
  service/           — бизнес-логика и валидация
migrations/          — SQL-миграции
.github/workflows/   — CI
```

## Разработка

Перед пушем:

```bash
make pre-push
```

Гоняет всё, что гоняет CI: линтер, юнит- и интеграционные тесты. Если локально зелёное — CI пройдёт.
```