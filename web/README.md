# AI Assistants Catalog — Frontend (SPA)

React-клиент каталога AI-ассистентов к backend-API из [`../api.yaml`](../api.yaml).

## Стек

| Задача | Решение | Почему |
| --- | --- | --- |
| UI | **React 18 + TypeScript** | TypeScript рекомендован как доп. задание; типы в `src/types/api.ts` зеркалят `api.yaml`, что ловит рассинхрон с бэком на этапе сборки. |
| Сборка | **Vite 5** | Быстрый dev-сервер и простой прокси на backend. |
| Роутинг | **react-router-dom 6** | Вложенные роуты + ленивые guard-компоненты по роли. |
| Серверное состояние | **TanStack Query 5** | Кэш, `isLoading/isError`, ретраи, инвалидация после мутаций «из коробки». |
| Клиентское/сессионное состояние | **Zustand** | Токен и роль — это один маленький глобальный кусок состояния; Zustand даёт минимум кода + `persist` в `localStorage`. Redux был бы избыточен. |
| HTTP | **Axios** | Интерцепторы для подстановки `Authorization` и централизованной обработки 401. |
| Тесты | **Vitest + Testing Library** | Быстрая замена Jest, работает на той же конфигурации Vite. |

## Запуск

### Локально (dev)

```bash
cd web
npm install
cp .env.example .env      # необязательно, дефолты подходят
npm run dev                # http://localhost:3000
```

Vite-дев-сервер проксирует `/api/*` на `http://localhost:8080` и срезает префикс
`/api` (см. `VITE_PROXY_TARGET`). Поэтому backend должен быть запущен на `8080`.
Отдельный префикс нужен, чтобы при обновлении страницы `/assistants`, `/runs/my`
или `/admin/...` dev-сервер отдавал SPA, а не проксировал навигацию в API.

### Проверки

```bash
npm run build      # tsc + vite build
npm run lint       # ESLint
npm test           # Vitest (10 тестов)
npm run coverage   # отчёт покрытия (v8)
```

### Docker

`web/Dockerfile` — multi-stage: сборка на `node:20-alpine`, раздача статики через `nginx`.
Nginx также проксирует API на контейнер `backend:8080` (см. `web/nginx.conf`),
поэтому `docker-compose.yaml` в корне репозитория должен содержать сервис `frontend`
(порт `3000:80`) и `backend`.

## Архитектура

```
src/
├── api/            # слой работы с API (единственное место, где собираются URL/заголовки)
│   ├── http.ts         # axios-инстанс, ApiError, интерцепторы, авто-разлогин на 401
│   ├── auth.ts         # POST /dummyLogin
│   ├── categories.ts   # GET/POST /categories
│   ├── assistants.ts   # CRUD /assistants + POST /assistants/{id}/run
│   ├── runs.ts         # GET /runs/my, GET /admin/runs
│   ├── queryKeys.ts    # ключи TanStack Query
│   └── index.ts        # wireHttpAuth() + реэкспорты
├── store/auth.ts   # Zustand + persist: token, user, role
├── hooks/queries.ts# useQuery/useMutation хуки по доменам
├── components/     # Layout, Guards, формы, таблицы, состояния, пагинация
├── pages/          # страницы роутера
├── lib/validation.ts
├── types/api.ts    # типы из api.yaml
└── router.tsx
```

### Хранение токена и роли

Единственный источник правды — **Zustand-стор `useAuthStore`** (`src/store/auth.ts`),
персистится в `localStorage` под ключом `ai-catalog-auth`.
Нигде больше JWT не дублируется. HTTP-слой получает актуальный токен лениво через
`setTokenGetter`, что исключает циклические импорты `api ↔ store`.
При ответе `401` интерцептор вызывает `clear()` — пользователь разлогинивается.

### Защита роутов

- `RequireAuth` — не пускает без токена, редирект на `/login` с сохранением `from`.
- `RequireRole roles={['admin']}` — оборачивает все `/admin/*`; обычный пользователь при
  ручном вводе URL будет перенаправлен на `/assistants`.

### Состояние фильтров

Фильтры каталога (`q`, `categoryId`, `includeInactive`, `page`) синхронизированы с
query-параметрами URL через `useSearchParams`, поэтому ссылку можно скопировать и
состояние восстановится. Поиск дебаунсится (~350 мс) перед записью в URL.

### Состояния экранов

На каждой странице с запросами есть `loading` (спиннер), `error` (с кнопкой «Повторить»)
и пустое состояние. Инвалидация: создание/редактирование ассистента и создание категории
сбрасывают кэш `['assistants']`/`['categories']`; запуск ассистента инвалидирует `['runs']`,
поэтому история обновляется автоматически.

## Скрытие системного промпта

`systemPrompt` — внутренняя настройка ассистента. **Обычному пользователю он не
показывается**: на странице `/assistants/:id` блок «Системный промпт» рендерится только
для роли `admin`. Backend по спецификации вправе не отдавать это поле обычному
пользователю (`systemPrompt` помечен как nullable). Так мы не раскрываем
системную инструкцию и не даём её обойти.

## Страницы

| Путь | Доступ | Назначение |
| --- | --- | --- |
| `/login` | все | выбор роли, вызов `/dummyLogin`, сохранение сессии |
| `/assistants` | авторизованные | каталог: поиск, фильтр по категории, пагинация, `includeInactive` для админа |
| `/assistants/:id` | авторизованные | карточка + форма запуска, вывод ответа/ошибки |
| `/runs/my` | авторизованные | мои запуски, фильтр по статусу, модалка с полным ответом |
| `/admin/runs` | admin | все запуски, фильтры по ассистенту и статусу |
| `/admin/categories/new` | admin | создание категории с валидацией |
| `/admin/assistants/new` | admin | создание ассистента |
| `/admin/assistants/:id/edit` | admin | редактирование + вкл/выкл |

## Тесты

10 тестов в 4 файлах (`Vitest` + `Testing Library`), сетевой слой замокан (`vi.mock('@/api')`),
компоненты и стор — настоящие:

- `pages/LoginPage.test.tsx` — успешный вход сохраняет токен и роль; ошибка входа показывается.
- `pages/AssistantsPage.test.tsx` — рендер карточек из API; пустое состояние; `includeInactive`
  доступен только админу.
- `pages/AssistantDetailPage.test.tsx` — системный промпт скрыт от user и виден admin;
  запуск ассистента отображает ответ.
- `components/AssistantForm.test.tsx` — валидация обязательных полей; успешный сабмит.

Покрытие (`npm run coverage`) — см. отчёт `coverage/`. Встроенный v8-репортер считает по
всем `src/**`; страницы, завязанные на реальные запросы и роутер целиком (`App`, `router`,
`AdminRunsPage`, `EditAssistantPage`, `NewCategoryPage`, `MyRunsPage`), в тестах напрямую
не монтируются, поэтому общий процент скромный — это осознанный компромисс ради стабильных
компонентных тестов. Ключевые пользовательские сценарии (вход, каталог, карточка+запуск,
форма создания) покрыты напрямую.

## Переменные окружения

| Переменная | По умолчанию | Назначение |
| --- | --- | --- |
| `VITE_API_BASE_URL` | `/api` | базовый URL API; по умолчанию относительный `/api` + прокси |
| `VITE_PROXY_TARGET` | `http://localhost:8080` | куда Vite проксирует API в dev-режиме |

В продакшн-контейнере API проксирует nginx (`web/nginx.conf`), поэтому браузер всегда
ходит на тот же origin, и CORS не нужен.
