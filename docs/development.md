# Разработка

## Инструменты и запуск

Версии: `mise.toml`, `go.mod`, `package.json`, `package-lock.json`. После `mise trust` и `mise install` запускайте команды через `mise exec --`.

`bin/configure github.com/owner/app` меняет Go module и импорты, имя npm-пакета и lock-файл. Источники skills и ссылки на исходный шаблон сохраняются. После настройки выполните `bin/ci`.

`bin/setup` копирует отсутствующий `.env`, устанавливает зависимости, поднимает PostgreSQL из `compose.yml`, применяет миграции и собирает frontend. Повторный запуск сохраняет `.env` и данные. Запускать Docker Compose следует из корня проекта.

`bin/dev` собирает Go и запускает его вместе с Vite. При остановке любого процесса завершается второй; Ctrl+C останавливает оба. React/CSS используют HMR. После правки Go нужен перезапуск. По умолчанию приложение — `http://localhost:3000`, Vite — `http://localhost:5173`. Порт Go можно изменить через `HTTP_ADDR` и `PUBLIC_URL` в `.env`; Vite берёт CORS origin из `PUBLIC_URL`. При изменении порта самого Vite обновите `vite.config.ts`, проверку hot origin в `internal/httpapp/app.go` и тесты вместе.

## Создать пользователя

```sh
mise exec -- bin/manage create-user --email you@example.com --admin
```

CLI читает одну строку пароля из stdin. Пароль должен занимать 12–72 байта; email нормализуется и уникален. Публичной регистрации и восстановления пароля нет.

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
| `DB_MAX_CONNECTIONS` | Положительное число соединений; default 10 |

`PUBLIC_URL` задаёт CORS origin Vite в development и проверяет deployment-конфигурацию; приложение не использует его как Host allowlist или генератор абсолютных ссылок. Ограничивайте публичные hostnames на reverse proxy. Локальный `.env` читают shell-обёртки; бинарники получают переменные из окружения и не загружают dotenv автоматически.

## SQL и миграции

Миграции в `internal/database/migrations` встраиваются в бинарник. Добавляйте следующий числовой SQL-файл, не меняйте уже применённые миграции. `mise exec -- bin/manage migrate` берёт advisory lock и применяет ожидающие файлы одной транзакцией. Ошибка откатывает изменения; сервер сам миграции не запускает. Нет автоматического down/rollback — обратное изменение оформляется новой миграцией.

Для локальной БД: `docker compose exec postgres psql -U starter -d starter_development`. Остановка без удаления данных: `docker compose stop postgres`.

## Компоненты

Стили и semantic tokens — `web/src/styles.css`; базовые компоненты — `web/src/components`. `npm run storybook` поднимает каталог на порту 6006, `npm run build:storybook` собирает его в `storybook-static`. Каталог не включается в production-контейнер. Storybook не создаёт hot file приложения.

Пакет shadcn не установлен. При необходимости примените `shadcn-inertia`, сохранив CSRF, реальные JSON-типы и CSP. Для обычных компонентов используйте исходные EM `sb-*` и `tailwind-best-practices`.
