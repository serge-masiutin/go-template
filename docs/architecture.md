# Архитектура

Go владеет URL, аутентификацией, авторизацией и данными. Gonertia передаёт имя React-страницы и props: первый запрос возвращает HTML, следующие Inertia-переходы — JSON. В production frontend и backend доступны с одного origin; отдельного публичного JSON API и React Router нет.

## Границы

- `cmd/server`: конфигурация, пул БД, session store, HTTP-сервер и завершение по сигналу.
- `cmd/worker`: River queues, отдельный LISTEN/NOTIFY connection, mail/AI adapters и graceful drain.
- `cmd/manage`: миграции и создание пользователя через тот же пакет `accounts`.
- `internal/accounts`, `internal/notes`: типизированные данные, валидация и GORM queries конкретной функции. Маленькому CRUD не нужны пустые service/repository/domain слои.
- `internal/httpapp`: маршруты, строгий декодер, CSRF, HTTP-авторизация, Inertia и session flash.
- `internal/database`: один pgx pool с SQL/GORM adapters; Goose и River migrations под deploy lock.
- `internal/notemail`, `internal/assistant`: явные операции, записи состояния, транзакционный enqueue и внешние adapters.
- `internal/background`: типизированные River workers, failure cleanup и периодическая сверка.
- `internal/observability`: Prometheus exporter с ограниченными labels.
- `web/src`: страницы, компоненты, DTO TypeScript и стили. `web/src/types.ts` — общие типы и единственное расширение `InertiaConfig`.

При появлении операции, которая координирует несколько записей и инвариантов, выделяйте именованный use case и явно задавайте транзакцию. Переносите проверки доступа в общий контракт операции, если она вызывается не только из HTTP. Внешние эффекты требуют отдельной гарантии доставки после commit; обычная goroutine её не даёт. River `InsertTx` сохраняет delivery intent вместе с данными операции. Подробности — [очередь и почта](background.md), [AI](ai.md).

## HTTP-контракт

| Маршрут | Доступ и поведение |
| --- | --- |
| `GET /login` | Форма входа |
| `POST /login` | Проверка пароля, смена session ID и CSRF, 303 на `/` |
| `POST /logout` | Уничтожение сессии, очистка Inertia history, 303 на `/login` |
| `GET /` | Вход обязателен; до 50 последних заметок текущего пользователя |
| `POST /notes` | 1–2000 Unicode-символов после trim; владелец берётся из сессии |
| `DELETE /notes/{id}` | Удаление только своей заметки; чужая/отсутствующая — 404 |
| `GET /admin` | Повторная проверка `admin` в БД; число пользователей |
| `GET /tools` | Свои почтовые и AI-операции, текущая доступность функций |
| `POST /tools/email` | Вход обязателен; письмо только себе, одна активная операция |
| `POST /tools/assistant` | Вход обязателен; 1–500 символов, одна активная операция |
| `GET /metrics` | Только при настроенном Bearer token, без session-аутентификации |
| `GET /health/live` | Процесс отвечает |
| `GET /health/ready` | PostgreSQL отвечает в пределах двух секунд |

Мутации принимают ограниченный JSON (16 KiB) и `X-CSRF-Token`. Неизвестные поля и лишний JSON отклоняются. Для файлов потребуется отдельный multipart-контракт и лимиты. Precognition-запросы отклоняются: validation-only обработчик не реализован.

Ошибки полей записываются в session flash, затем следует 303 и обычный GET с `errors`. Поля ошибок плоские; named error bags автоматически не поддерживаются. Неверный JSON — 400, отсутствие права — 403/404, неожиданный сбой — 500. `net/http.CrossOriginProtection` дополняет session CSRF-проверку.

Публичные ID сериализуются строками. Password hash, cookie и приватные конфигурационные значения в props не попадают. CSRF передаётся через `Always`, включая partial reload. Пользователь — prop конкретной страницы, не глобальная переменная. Каждый повторный запрос заново проверяет доступ.

## Сессии и эксплуатация

SCS хранит сессии в PostgreSQL: 24 часа абсолютного срока и 2 часа бездействия. В production cookie имеет Secure; HttpOnly и SameSite=Lax действуют во всех средах. Права администратора читаются заново на каждом защищённом запросе.

Вход ограничен десятью попытками в минуту на прямой адрес соединения; таблица ограничена 4096 адресами. Forwarded-заголовки не считаются доверенным источником IP. За reverse proxy адрес может быть общим: настройте ограничение на внешнем входе и продумайте общий лимит при нескольких экземплярах приложения.

Production игнорирует Vite hot file, использует manifest и версию ресурсов. При несовпадении версии Gonertia отвечает 409 для полной перезагрузки. CSP разрешает same-origin scripts; inline JavaScript по умолчанию запрещён. Ошибки логируются безопасными категориями и SQLSTATE, без сырых driver messages с пользовательскими значениями. Panic на HTTP-границе даёт 500 и stack trace без значения panic.

## Как развивать

`layered-go` сохраняет методику оригинального `layered-rails`, но заменяет framework-механизмы на Go-контракты. Четыре концептуальных слоя не требуют четырёх деревьев каталогов. Для Inertia-функций начните с `inertia-go-architecture`, затем загрузите профильный skill. Поддержку deferred, uploads, shadcn, новых каналов уведомлений и других возможностей добавляйте вместе с обработкой ошибок и тестами; их примеры не означают, что они установлены.


## Persistence и фоновые эффекты

GORM records и публичные DTO отделены: password hash не является полем `accounts.User`, для фоновых операций выбираются только публичные поля, которые явно переводятся в DTO. Не принимайте ORM records напрямую из HTTP. GORM query fragments, таблицы и sort expressions принадлежат серверу; пользовательские значения передаются bind-параметрами.

GORM `Transaction` предоставляет общий `*sql.Tx` для нескольких записей и River enqueue. Внешний SMTP/model call происходит после commit и после atomic claim состояния. Шаблон не добавляет business callbacks, AutoMigrate или lazy associations. Pure domain behavior не зависит от GORM/HTTP.
