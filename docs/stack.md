# Набор библиотек

Проверено 27 сентября 2026 года по документации разработчиков и опубликованным Go-модулям. Точные версии приложения закреплены в `go.mod`, CLI-инструментов — в `tools/go.mod`, контейнеров — в `compose.yml`. Это выбор для данного Go + Inertia + React starter, а не универсальный рейтинг библиотек.

## Что уже подключено

| Возможность Rails-стека | Выбор в Go | Где используется |
| --- | --- | --- |
| Rails routes/controllers + Inertia | `net/http` + Gonertia 3.0.0 | Маршруты, middleware, страницы, формы и session flash |
| Active Record | GORM 1.31.2, PostgreSQL driver 1.6.3, pgx 5.11.0 | Accounts, заметки, почтовые операции и AI-запросы; generics API и явные projections |
| Миграции | Goose 3.28.0 | SQL-файлы, транзакция на миграцию, общий deploy lock с миграциями River |
| Solid Queue / Active Job | River 0.47.0 | Транзакционная постановка mail/AI, worker, периодическая сверка состояния |
| Mission Control Jobs | River UI 0.19.0 OSS | Локальный интерфейс очередей через Compose profile `tools` |
| Action Mailer / Active Delivery | go-mail 0.8.1 + `notemail.Service` | SMTP, TLS, текстовый шаблон и явная доставка через River |
| Локальный просмотр писем | Mailpit 1.31.2 | SMTP sink и web UI; исходящие письма не пересылаются в интернет |
| RubyLLM / Active Agent | Genkit Go 1.13.1 | Типизированный flow, JSON output, ограниченный цикл tools; Gemini и OpenAI |
| Anyway Config | caarlos0/env 11.4.1 | Типизированные nested structs, defaults и отдельная проверка инвариантов |
| Sessions | SCS 2.9.0 + pgxstore | Серверные сессии в том же PostgreSQL pool |
| Yabeda / Prometheus | prometheus/client_golang 1.24.1 + httpsnoop | HTTP latency/status, runtime, соединения БД, состояния очередей |
| Reload в development | Air 1.67.4 + Vite | Пересборка Go server/worker, React HMR, совместная остановка |
| Аудит зависимостей | govulncheck 1.8.0 + npm audit | Go-анализ достижимых уязвимых функций и npm-зависимости |

Пользовательский сценарий — `/tools`: отправить свои последние заметки на свой email либо задать AI вопрос по этим заметкам. Конфигурация провайдера остаётся на сервере. В очередь попадают только ID операций; вопрос и ответ принадлежат пользователю и не выдаются другим аккаунтам.

## Почему этот набор

**HTTP.** Современный [ServeMux поддерживает методы и параметры маршрутов](https://go.dev/blog/routing-enhancements). Gonertia уже связывает сервер с React. [Chi совместим с net/http](https://github.com/go-chi/chi) и подходит при потребности в более сложных группах маршрутов; второй HTTP framework сейчас не даёт этому приложению нужной дополнительной возможности. Можно применять обычные net/http middleware без переписывания handlers.

**ORM.** [GORM рекомендует generics API](https://gorm.io/docs/the_generics_way.html): операции возвращают типизированный результат и error, неоднозначные `Save` и `FirstOrCreate` исключены из этого API. [Пользовательские значения передаются placeholders](https://gorm.io/docs/security.html); имена колонок, таблиц и сортировки должны оставаться серверными. В starter нет AutoMigrate, persistence callbacks с бизнес-эффектами, auto-save associations или сериализации ORM-record целиком.

[Ent](https://entgo.io/docs/getting-started/) — сильная альтернатива для проектов, которым нужен сгенерированный schema-first API и большой граф связей. Здесь выбран меньший объём генерации и обычные SQL-миграции. [Bun](https://bun.uptrace.dev/guide/queries.html) удобен для SQL-first запросов, но [его placeholders интерполируют и экранируют значения, в том числе удаляют NUL](https://bun.uptrace.dev/guide/placeholders.html); для этого шаблона выбран bind-подход GORM/pgx и громкая ошибка на некорректных данных.

**Очередь.** [River поддерживает общую транзакцию с GORM](https://riverqueue.com/docs/gorm). Запись операции и job либо сохраняются вместе, либо вместе откатываются. Это закрывает описанное [Brandur Leach окно потери работы между commit и enqueue](https://brandur.org/job-drain). Отдельный Redis не требуется. У SMTP и модельного API нет общей транзакции с PostgreSQL: поэтому starter не обещает exactly-once внешних эффектов и не повторяет неопределённую отправку автоматически.

**AI.** [Genkit tools](https://genkit.dev/docs/go/tool-calling/) дают проверяемые типы и готовый цикл вызовов. Используется стабильный API flows/generate, без экспериментального `genkit/exp`: [full-stack Agents API пока beta](https://genkit.dev/docs/go/agents/overview/). [Google ADK](https://github.com/google/adk-go) и [Eino](https://www.cloudwego.io/docs/eino/overview/) дают более широкий orchestration; он не нужен одному помощнику с read-only инструментом. Поддержка инструментов, ограничения модели и тесты важнее количества agent frameworks в одном проекте.

**Наблюдаемость.** Genkit может писать полный контент в telemetry. В приложении exporter и reflection UI не включены; переменные, автоматически открывающие этот канал, отклоняются. Это следует из [описанного Genkit разделения traces и content logs](https://genkit.dev/docs/go/observability/telemetry-collection/). Метрики и логи приложения содержат безопасные операционные поля. Письма и ответы можно проверить через пользовательский интерфейс и локальный Mailpit.

## Расширение

Authorization пока выражается owner-scoped запросами и проверкой актуального admin-флага. Casbin/OpenFGA нужны при реальной модели policy/relationship access, а не для замены двух явных проверок. Realtime, S3/uploads, платежи, vector search, регистрация и восстановление пароля остаются отдельными продуктовыми задачами. Не включайте альтернативные ORM, agent runtimes и очереди одновременно без конкретного контракта между ними.
