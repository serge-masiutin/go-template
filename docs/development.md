# Разработка

## Инструменты и запуск

Версии: `mise.toml`, `go.mod`, `package.json`, `package-lock.json`. После `mise trust` и `mise install` запускайте команды через `mise exec --`.

`bin/configure github.com/owner/app` меняет Go modules приложения/инструментов и импорты, имя npm-пакета и lock-файл. Источники skills и ссылки на исходный шаблон сохраняются. После настройки выполните `bin/ci`.

`bin/setup` копирует отсутствующий `.env`, устанавливает зависимости, поднимает PostgreSQL и Mailpit из `compose.yml`, применяет миграции и собирает frontend. Повторный запуск сохраняет `.env` и данные. Запускать Docker Compose следует из корня проекта.

`bin/dev` запускает Vite и два Air watcher: server и worker. Go пересобирается при изменении связанных файлов; React/CSS используют HMR. Ctrl+C завершает все процессы; аварийное завершение watcher/Vite останавливает остальных. Ошибка компиляции видна в терминале, Air ждёт исправления. Конфигурация `.env` читается при запуске supervisor: после её изменения перезапустите `bin/dev`. По умолчанию приложение — `http://localhost:3000`, Vite — `http://localhost:5173`. Порт Go можно изменить через `HTTP_ADDR` и `PUBLIC_URL` в `.env`; Vite берёт CORS origin из `PUBLIC_URL`. При изменении порта самого Vite обновите `vite.config.ts`, проверку hot origin в `internal/httpapp/app.go` и тесты вместе.

## Создать пользователя

```sh
mise exec -- bin/manage create-user --email you@example.com --admin
```

CLI читает одну строку пароля из stdin. Пароль должен занимать 12–72 байта; email нормализуется, уникален и ограничен 254 байтами. Публичной регистрации и восстановления пароля нет.

Чтобы скрыть ввод в терминале, выполните через Bash:

```sh
bash -c 'read -r -s -p "Password: " starter_password; printf "\n" >&2; printf "%s\n" "$starter_password" | mise exec -- bin/manage create-user --email you@example.com --admin; unset starter_password'
```

Не храните production-пароль в shell history, аргументах команды или файле репозитория. Уберите `--admin`, чтобы создать обычного пользователя.

## Конфигурация

| Переменная | Контракт |
| --- | --- |
| `APP_ENV` | `development`, `test`, `production`; пустое значение — development |
| `DATABASE_URL` | Обязательная строка подключения PostgreSQL |
| `HTTP_ADDR` | Адрес listener; default `127.0.0.1:3000` |
| `PUBLIC_URL` | HTTP origin без пути; production требует явного HTTPS origin |
| `DB_MAX_CONNECTIONS` | Не менее 2 соединений; default 10; server и worker имеют отдельные бюджеты |
| `WORKER_CONCURRENCY` | 1–100 mail workers; default 4; AI — один на процесс |
| `SHUTDOWN_TIMEOUT` | 1s–1m; default 10s |
| `METRICS_TOKEN` | Пустое значение выключает `/metrics`; иначе минимум 32 символа |

`PUBLIC_URL` задаёт CORS origin Vite в development и проверяет deployment-конфигурацию; приложение не использует его как Host allowlist или генератор абсолютных ссылок. Ограничивайте публичные hostnames на reverse proxy. Локальный `.env` читают shell-обёртки; бинарники получают переменные из окружения и не загружают dotenv автоматически.

## SQL и миграции

Goose читает SQL-файлы `internal/database/migrations`, встроенные в бинарник. Добавляйте следующий числовой файл с `-- +goose Up` / `-- +goose Down`, не меняйте применённые миграции. `mise exec -- bin/manage migrate` держит один advisory lock на миграции приложения и River в отдельной PostgreSQL-сессии; она закрывается после миграций и не возвращается в application pool. Каждая миграция приложения атомарна; весь набор не является одной транзакцией. При ошибке исправьте причину и повторите команду; уже применённые версии сохраняются.

Первая Go-миграция принимает известную историю прежнего starter либо создаёт исходные таблицы; неизвестная legacy history вызывает ошибку. Исходный DDL сохранён в `internal/database/bootstrap`. Сервер/worker миграции не запускают. Management CLI не предоставляет destructive down: обратное изменение оформляйте новой миграцией.

GORM использует generics API и тот же pgx pool через `stdlib.OpenDBFromPool`. Запросы конкретной функции остаются в её store/operation. Не включайте AutoMigrate, SQL debug logging или business hooks. Многозаписочные операции открывают явную транзакцию; пример совместного `InsertTx` — `internal/notemail/notemail.go`.

Для локальной БД: `docker compose exec postgres psql -U starter -d starter_development`. Остановка без удаления данных: `docker compose stop postgres`.

## Компоненты

Стили и semantic tokens — `web/src/styles.css`; базовые компоненты — `web/src/components`. Tailwind сканирует только `web/src`, чтобы документация и skills не влияли на CSS. `npm run storybook` поднимает каталог на порту 6006, `npm run build:storybook` собирает его в `storybook-static`. Каталог не включается в production-контейнер. Storybook не создаёт hot file приложения.

Пакет shadcn не установлен. При необходимости примените `shadcn-inertia`, сохранив CSRF, реальные JSON-типы и CSP. Для обычных компонентов используйте исходные EM `sb-*` и `tailwind-best-practices`.


## Дополнительные инструменты

- Почта и River UI: [очередь и почта](background.md); Mailpit — `http://localhost:8025`, River UI — `http://localhost:8087` после включения Compose profile.
- AI-конфигурация и тесты: [AI-помощник](ai.md).
- `bin/tool air -v` и `bin/tool govulncheck ./...` используют отдельный `tools/go.mod`. Это позволяет обновлять dev tooling без изменения runtime dependencies. После обновления tools выполните `go -C tools mod tidy`.
- `/metrics` требует `Authorization: Bearer <METRICS_TOKEN>`. Prometheus exporter включает Go runtime/process, пул БД, HTTP duration/status по шаблону маршрута и состояния River jobs. Он не экспортирует тела запросов, notes, вопросы, ответы или сырые URL. Ошибка чтения queue metrics делает scrape неуспешным, а не выдаёт нулевую очередь.

Установка Prometheus/Grafana и OTLP collector в starter не входит. Для мониторинга подключите exporter к существующей системе и отдельно контролируйте доступность worker/возраст очереди.
