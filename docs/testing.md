# Проверки

## Основной набор

```sh
mise exec -- bin/ci
```

Включает форматирование Go, `go mod tidy -diff`, `go mod verify` для обоих модулей, `go vet`, `go test -race`, TypeScript, Vite build, Vitest, Storybook build и проверку skills. PostgreSQL для этого набора не нужен. Исходники frontend форматируются Prettier; проверка входит в CI.

## PostgreSQL и браузер

Поднимите локальную БД через `bin/setup`, затем один раз создайте отдельную тестовую:

```sh
docker compose exec postgres createdb -U starter starter_test
```

Если она уже существует, используйте её. Тестовые команды откажутся работать с именем БД без суффикса `_test`. Каждый запуск создаёт собственную схему и удаляет её после завершения; development-схема не затрагивается.

```sh
export TEST_DATABASE_URL='postgres://starter:local-development-only@127.0.0.1:5437/starter_test?sslmode=disable'
mise exec -- go test -race -tags=integration ./...
mise exec -- npx playwright install chromium
mise exec -- bin/test-browser
```

Предварительно выполните `bin/ci`, чтобы иметь актуальный manifest. `bin/test-browser` создаёт изолированную схему и синтетического администратора, запускает Playwright и удаляет схему. Порт 3100 должен быть свободен. `APP_ENV=test` использует собранные assets и игнорирует development hot file. В Linux для браузера могут потребоваться системные зависимости: `npx playwright install --with-deps chromium`.

Integration tests проверяют миграции, вход, одноразовые ошибки форм, CSRF/его ротацию, межсайтовые запросы, строгий JSON, version mismatch, partial props, разделение заметок по владельцу, отзыв admin-права, logout и отсутствие мутаций при неподдерживаемой Precognition-валидации. Браузер проходит неудачный/успешный вход, валидацию заметки, создание, перезагрузку, удаление и выход; JS-ошибки приводят к падению.

## Изменение контрактов

Тестируйте публичное поведение и негативные сценарии. SQL/session/authorization проверяйте с настоящим PostgreSQL. Ошибки типов и браузерный сценарий дополняют Go-тесты, а не заменяют их. Для новой deferred-функции добавьте отказ загрузки, retry и повторную авторизацию; для файлов — фактический multipart payload и лимиты.

GitHub Actions выполняет оба набора и собирает Docker image. Отдельный `bin/tool govulncheck ./...` также выполняется в GitHub Actions. В основной набор не входят нагрузочные испытания, production smoke на вашем хостинге и платные генерации реальных AI-провайдеров.


## Очередь, почта и AI

`go test ./internal/assistant ./internal/mailing` использует синтетические данные, локальные HTTP/SMTP-серверы и фальшивую модель для проверки контрактов. Проверяются tool loop limit, схема ответа, отсутствие ответа без чтения заметок, Gemini thought signature, auth headers, token budget и отсутствие HTTP retry при 503. Никакие реальные письма или платные генерации не выполняются.

PostgreSQL-набор дополнительно проверяет принятие legacy migration history, rollback записи при ошибке enqueue, конкурентную постановку, owner scope, двойной запуск, удаление аккаунта, работу настоящего River worker, failure cleanup и reconciliation после аварии. Метрики проверяются на Bearer-доступ и отсутствие пользовательского контента.

Качество ответов реальной модели — отдельная [оценка](../evals/notes-assistant.md). HTTP fixtures не доказывают качество модели. `govulncheck` может отмечать неиспользуемые подпакеты внутри модулей: оценивайте отдельно reachable symbols, imported packages и module-only findings; не называйте их одним числом.
