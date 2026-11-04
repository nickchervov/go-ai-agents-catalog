# AI Assistants Catalog

Fullstack-каталог AI-ассистентов: администратор управляет категориями и ассистентами, пользователь просматривает каталог, фильтрует по категории, ищет по названию и описанию и запускает ассистента с пользовательским контекстом.

**Frontend SPA написан с помощью ИИ.** Основная задача проекта — backend: REST API по [`api.yaml`](api.yaml), JWT-авторизация, PostgreSQL, миграции, мок-LLM и история запусков. Frontend использовался как клиент для проверки контракта.

## Стек

| Слой | Технологии |
| --- | --- |
| Backend | Go 1.26, chi, pgx/v5 + sqlx, golang-migrate, JWT (HS256), slog |
| БД | PostgreSQL 16 |
| LLM | встроенный детерминированный mock; опционально OpenAI-compatible API |
| Frontend | React 18, TypeScript, Vite 5, react-router-dom 6, TanStack Query 5, Zustand, Axios |
| Инфраструктура | Docker Compose, multi-stage Dockerfile, healthchecks |

## Быстрый старт

Требуется Docker Compose v2 (или `docker-compose` v1).

```bash
docker compose up --build
# либо:
docker-compose up --build
```

После запуска:

- frontend — <http://localhost:3000>
- backend — <http://localhost:8080>
- healthcheck — <http://localhost:8080/_info> (200 OK)

Миграции применяются автоматически при старте backend. Данные хранятся в volume `pgdata`. Для чистой проверки схемы:

```bash
docker compose down -v
docker compose up --build
```

Для локального запуска backend:

```bash
cp .env.example .env   # заполнить DSN, при необходимости LLMAPI/BASE_URL
go run ./cmd/api
```

Для локального frontend:

```bash
cd web
npm install
npm run dev            # http://localhost:3000
```

Vite проксирует `/api/*` на `http://localhost:8080` и срезает префикс. В production то же делает nginx.

## Переменные окружения

Все значения по умолчанию заданы в `docker-compose.yaml`.

| Переменная | Назначение | Значение по умолчанию |
| --- | --- | --- |
| `PORT` | порт backend | `:8080` |
| `DSN` | подключение к PostgreSQL | `postgres://postgres:1234@db:5432/agent_catalog` |
| `LLMAPI` | API-ключ LLM; пусто — mock | пусто |
| `BASE_URL` | base URL OpenAI-compatible API; пусто — mock | пусто |
| `LOGGER_FORMAT` | формат логов (`json`/`text`) | `json` |
| `SECRET` | секрет подписи JWT | значение в compose |
| `TTL_TOKEN_MIN` | время жизни JWT, минуты | `60` |

Frontend:

| Переменная | Назначение | Значение по умолчанию |
| --- | --- | --- |
| `VITE_API_BASE_URL` | базовый URL API | `/api` |
| `VITE_PROXY_TARGET` | куда Vite проксирует `/api` | `http://localhost:8080` |

`.env` не коммитится; в репозитории есть `.env.example`.

## API

Контракт — [`api.yaml`](api.yaml). Авторизация — `Authorization: Bearer <JWT>`. Ошибки возвращаются в конверте:

```json
{"error":{"code":"INVALID_REQUEST","message":"..."}}
```

| Метод | Путь | Доступ | Назначение |
| --- | --- | --- | --- |
| POST | `/dummyLogin` | все | JWT по роли `admin`/`user` |
| GET | `/categories` | авторизованные | список категорий |
| POST | `/categories` | admin | создание категории |
| GET | `/assistants` | авторизованные | каталог: `categoryId`, `q`, `includeInactive`, `page`, `pageSize` |
| POST | `/assistants` | admin | создание ассистента |
| GET | `/assistants/{assistantId}` | авторизованные | карточка ассистента |
| PUT | `/assistants/{assistantId}` | admin | редактирование; `isActive` обязателен |
| POST | `/assistants/{assistantId}/run` | авторизованные | запуск; в теле только `userPrompt` |
| GET | `/runs/my` | авторизованные | история запусков текущего пользователя |
| GET | `/admin/runs` | admin | все запуски |
| GET | `/_info` | все | healthcheck |

Параметры пагинации: `/assistants` — `pageSize=10` по умолчанию, `/runs/my` и `/admin/runs` — `20`; максимум `100`. Для обычного пользователя `includeInactive` игнорируется. По умолчанию каталог показывает только активных ассистентов.

## Авторизация

- JWT HS256, claims: `user_id`, `email`, `role`, `exp`.
- `/dummyLogin` возвращает фиксированные UUID: `admin` — `11111111-1111-1111-1111-111111111111`, `user` — `22222222-2222-2222-2222-222222222222`.
- Проверка токена и ролей — в middleware; admin-роуты закрыты централизованно.
- При запуске ассистента `user_id` берётся из токена, а не из тела запроса.
- `systemPrompt` — внутренняя настройка: администратор получает его в API, обычный пользователь — пустое значение. Frontend его не отображает.

## Запуск ассистента и история

1. Проверяется, что `userPrompt` не пустой после `TrimSpace`.
2. Загружается ассистент; если не найден — 404, если выключен — 409.
3. В `assistant_run` создаётся запись со статусом `pending`.
4. Вызывается LLM с системным промптом, моделью и пользовательским промптом.
5. Успех — `success` и `result`; ошибка, таймаут или пустой ответ — `failed` и `error`.
6. Ответ API — 201 с данными запуска. При ошибке провайдера API возвращает 502, но запуск остаётся в истории.

Запись обновляется по `id` (upsert). Каждый запуск сохраняется независимо от результата LLM. Если процесс упадёт между `pending` и финальным обновлением, запуск останется в `pending` — для тестового задания это допустимо.

## LLM: mock и реальный провайдер

По умолчанию используется `mocks.MockLLM`: без сети, детерминированный ответ `"Nice prompt take result back"`. Это позволяет запускать проект и тесты без ключей.

Контракт провайдера:

```go
SendMessage(ctx context.Context, model, systemPrompt, userPrompt string) (string, error)
```

Backend формирует запрос из трёх полей ассистента и пользовательского ввода. Пример итогового запроса к OpenAI-compatible API:

```json
{
  "model": "gpt-4o-mini",
  "instructions": "Ты — повар. Составь простой домашний рецепт из ингредиентов пользователя: название, время и 3–5 шагов.",
  "input": "курица, лаваш, огурцы, томаты, соус"
}
```

`instructions` — `systemPrompt`, `input` — `userPrompt`. Frontend не передаёт системный промпт.

Переключение на реальный провайдер: в `docker-compose.yaml` задать непустые `LLMAPI` и `BASE_URL` (или подключить override-файл):

```yaml
# docker-compose.override.yaml
services:
  backend:
    environment:
      LLMAPI: "sk-..."
      BASE_URL: "https://api.openai.com/v1"
```

```bash
docker compose -f docker-compose.yaml -f docker-compose.override.yaml up --build
```

Если `LLMAPI` или `BASE_URL` пусты, backend принудительно использует mock. Поддерживается любой OpenAI-compatible endpoint, включая локальный.

Обработка сбоев:

- вызов LLM ограничен `context.WithTimeout` — 30 секунд;
- ошибка провайдера → запуск `failed`, ответ 502 `LLM_PROVIDER_ERROR`;
- пустой или пробельный ответ → запуск `failed`, ответ 502;
- ошибка БД при сохранении → 500 `INTERNAL_ERROR`.

## База данных

Миграции — в `pkg/migrations`, применяются автоматически при старте. Схема:

- `users` — `id`, `email` (unique), `role` (`admin`/`user`), `created_at`;
- `categories` — `id`, `category` (unique), `description`, `created_at`;
- `assistants` — `id`, `name`, `description`, `category_id` (FK → `categories`, `NOT NULL`), `model`, `system_prompt` (`NOT NULL`), `ex_user_prompt`, `is_activated`, `created_at`, `updated_at`;
- `assistant_run` — `id`, `assistant_id` (FK, `NOT NULL`), `user_id` (FK, `NOT NULL`), `user_prompt`, `result`, `status` (`pending`/`success`/`failed`), `error`, `created_at`.

Индексы:

- `idx_user_role` — по роли пользователя;
- `idx_assistants_category_id` — фильтр каталога по категории;
- `idx_assistants_search_trgm` — pg_trgm для поиска по названию/описанию;
- `idx_assistant_run_status`, `idx_assistant_run_user_id`, `idx_assistant_run_created_at`, `idx_assistant_run_assistant_id` — история запусков и фильтры.

`total` для пагинации считается отдельным `COUNT(*)`, поэтому корректно возвращается и на странице за пределами выдачи. Внешние ключи защищают ассистента от несуществующей категории, а запуск — от несуществующего ассистента/пользователя.

## Frontend

Страницы: `/login`, `/assistants`, `/assistants/:id`, `/runs/my`, `/admin/runs`, `/admin/categories/new`, `/admin/assistants/new`, `/admin/assistants/:id/edit`; неизвестный путь — 404.

Решения:

- серверное состояние — TanStack Query: кэш, loading/error, инвалидация после создания, редактирования и запуска;
- токен и роль — Zustand + `persist` в `localStorage` под ключом `ai-catalog-auth`; это единственное место хранения JWT;
- Axios-интерцептор подставляет `Authorization` и при 401 очищает сессию;
- admin-роуты защищены guard-компонентом;
- фильтры каталога синхронизированы с query-параметрами URL;
- на страницах с запросами есть loading/error/empty состояния;
- системный промпт не показывается обычному пользователю.

Проверки:

```bash
cd web
npm test          # 10 тестов
npm run coverage  # покрытие v8
npm run build
npm run lint
```

Текущее покрытие frontend: **36.77% statements, 69.93% branches, 43.51% functions, 36.77% lines**. Основные непокрытые зоны — admin-страницы, таблица запусков и layout; их нужно добавить в тесты.

## Тесты

Команды:

```bash
gofmt -l .
go vet ./...
go test ./... -cover
cd web && npm test && npm run coverage
```

Обязательный минимум по заданию:

**Backend, unit:**

- `RunAssistant`: успешный запуск; ошибка LLM; пустой ответ; неактивный ассистент; пустой `userPrompt`; ассистент не найден;
- `CreateAssistant`: `isActive` по умолчанию `true`; валидация обязательных полей; несуществующая категория;
- `UpdateAssistant`: отсутствующий `isActive` → 400; успешное обновление; несуществующая категория; ассистент не найден;
- `Assistants`: фильтр по `categoryId`; `includeInactive` только для admin; поиск; пагинация и `total`;
- `DummyLogin`: фиксированные UUID и роли; некорректная роль.

**Backend, интеграционный/E2E:**

- `admin создал категорию → admin создал ассистента → user запустил ассистента → запуск есть в истории`;
- фильтрация ассистентов по `categoryId` возвращает только выбранную категорию;
- запуск неактивного ассистента запрещён (409);
- негативные сценарии: 401 без токена, 403 для пользователя на admin-роутах, ошибка LLM.

E2E следует запускать на отдельной тестовой БД, а не на основной.

**Frontend:** минимум 3 теста на компоненты/сценарии — логин, каталог и фильтры, карточка/запуск, валидация формы. Сейчас в репозитории 10 тестов в 4 файлах.

Coverage backend добавляется вместе с тестами: `go test ./... -cover`. После реализации тестов значения coverage нужно обновить в этом разделе.

## Принятые решения

- Домен (`internal/domain`) не содержит `db`-тегов; row-модели и маппинг — в `internal/adapters/postgres/models.go`; HTTP-DTO — в `internal/dto`.
- Репозиторий возвращает `total` отдельно, а не через `COUNT(*) OVER()`.
- `isActive` при создании по умолчанию `true`; при обновлении поле обязательно, иначе 400.
- Mock-LLM вынесен в `mocks/`, чтобы production-код не зависел от тестовой реализации.
- Один `AgentService` и один репозиторий — осознанный выбор для размера тестового задания; дробление на микросервисы не требуется.
- Frontend написан с помощью ИИ; основные усилия направлены на backend, контракт и БД.
