# Go Template: контракт работы агента

- Отвечай и пиши документацию по-русски. Код, технические сообщения, комментарии и skills — по-английски; `clear-writing` — по-русски.
- Сначала прочитай связанные файлы, `go.mod`, `package-lock.json` и профильный документ. Меняй контракт вместе со всеми потребителями.
- Стек: Go `net/http`, Gonertia v3, GORM generics/pgx, Goose, SCS, River, go-mail, Genkit; React, TypeScript, Vite, Tailwind, Storybook. Сервер владеет маршрутами, доступом и данными; React — представлением и локальными взаимодействиями.
- Используй доменные пакеты `internal`; не создавай по сервису, интерфейсу и репозиторию на каждый endpoint. См. `layered-go` и `docs/architecture.md`.
- На HTTP-границе декодируй строгий ограниченный JSON. Не принимай владельца из тела вместо аутентифицированного пользователя. SQL параметризуй; транзакционные запросы выполняй через тот же `*sql.Tx`/GORM transaction; River enqueue — `InsertTx` внутри неё.
- Не включай AutoMigrate, SQL debug logs или business effects в ORM hooks. Mail/AI идут через River; payload содержит ID, worker повторно устанавливает scope, внешние вызовы выполняются вне транзакции.
- Каждая мутация передаёт `X-CSRF-Token`. Ошибка формы — session errors + 303; неверный JSON — 400; отсутствие доступа — 403/404. Precognition не поддерживается и отклоняется.
- Авторизуй каждый запрос, включая partial/deferred. Не клади пользователя в глобальные props; используй request context. Loaders Gonertia могут исполняться параллельно: pool допустим, общий Tx/Conn — нет.
- DTO перечисляют публичные поля. ID в JSON — строки. Обновляй Go JSON и TypeScript вместе. Автоматической генерации типов нет.
- Ошибки не подменяй пустыми результатами. Логируй безопасный контекст и категории; не записывай пароли, DSN, session/CSRF tokens или пользовательские тексты.
- Не переписывай оригинальные EM skills под проект. Их файлы проверяются по SHA-256. Проектные ограничения — здесь; адаптации и источники — `docs/skills.md`.

## Контекст

| Задача | Документ | Skills |
| --- | --- | --- |
| Архитектура Go | `docs/architecture.md` | `layered-go` |
| Библиотеки, ORM, очередь, AI | `docs/stack.md`, `docs/background.md`, `docs/ai.md` | `layered-go`, его `references/installed-stack.md` |
| Страницы, маршруты, формы | `docs/architecture.md` | `inertia-go-architecture`, затем соответствующий `inertia-go-*` |
| JSON и TypeScript | `docs/architecture.md` | `go-serialization`, `inertia-go-typescript` |
| Компоненты и Storybook | `docs/development.md` | `tailwind-best-practices`, `sb-hub`, затем нужный `sb-*` |
| Проверки | `docs/testing.md` | `inertia-go-testing` |
| Запуск и инфраструктура | `docs/development.md`, `docs/deployment.md` | `inertia-go-setup` |
| Медленный запуск | `docs/development.md` | `go-boot-profiling` |
| README и текст | `README.md` | `good-readme`, `clear-writing` |

## Проверки

Команды запускай через `mise exec --`. Для кода и конфигурации выполняй `bin/ci`; для HTTP/SQL/session-контрактов — также integration tests и `bin/test-browser` по `docs/testing.md`. Проверяй изменённое поведение и негативные сценарии. Сообщай, что выполнено и что осталось непроверенным. Не завершай работу со stubs или замаскированными ошибками.
