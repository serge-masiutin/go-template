# Контейнер и деплой

Образ собирает frontend и три Go-бинарника, затем запускает сервер от nonroot в distroless. Node, компилятор, Storybook и skills в runtime image не попадают.

```sh
docker build -t my-app .
```

Задайте `DATABASE_URL`, `PUBLIC_URL=https://app.example.com`; по умолчанию образ использует `APP_ENV=production` и `HTTP_ADDR=0.0.0.0:3000`. Секреты передавайте средствами платформы. `compose.yml` предназначен для локальной разработки: его пароль и отключённый TLS не подходят для публичной БД.

Пример формы команд, где переменные уже заданы безопасным способом:

```sh
docker run --rm --env DATABASE_URL --env PUBLIC_URL --entrypoint /app/manage my-app migrate
docker run --rm -i --env DATABASE_URL --env PUBLIC_URL --entrypoint /app/manage my-app create-user --email you@example.com --admin
docker run --rm -p 127.0.0.1:3000:3000 --env DATABASE_URL --env PUBLIC_URL my-app
docker run --rm --env DATABASE_URL --env PUBLIC_URL --entrypoint /app/worker my-app
```

Команда создания пользователя читает пароль из stdin. Выполняйте миграции отдельным шагом перед запуском новой версии. Для меняющегося приложения проектируйте миграции совместимыми с одновременно работающими версиями; наличие транзакции не делает любое DDL-изменение безопасным для rolling deploy.

Reverse proxy завершает TLS, сохраняет публичный Host, ограничивает допустимые hostnames и размер/скорость запросов. Secure cookie требует HTTPS в браузере. Не открывайте базу, профилировщик и внутренние интерфейсы в интернет. Встроенный login limiter локален процессу и видит адрес соединения; настройте внешний лимит с учётом proxy и числа экземпляров.

Проверки состояния: `/health/live` проверяет процесс, `/health/ready` — подключение PostgreSQL. SIGTERM запускает graceful shutdown с пределом `SHUTDOWN_TIMEOUT` (по умолчанию 10 секунд). Worker дополнительно получает 5 секунд для отмены после drain timeout; дайте платформе больший termination grace period. Модельный и SMTP timeout ограничены отдельно. Делайте резервные копии PostgreSQL и проверяйте восстановление.

Шаблон не создаёт инфраструктуру и не публикует приложение на хостинг. Для реального окружения отдельно проверьте DNS/TLS, секреты, права БД, proxy, backups и поведение обновления. Production CSP блокирует inline scripts; новые внешние сервисы требуют осознанного изменения политики и браузерной проверки.


Server и worker должны получать согласованные MAIL/AI settings, описанные в [почте](background.md) и [AI](ai.md); добавьте соответствующие `--env`/secret bindings к примеру выше. Базовые команды оставляют обе функции выключенными. Runtime pool budget задаётся на процесс; worker дополнительно открывает одно LISTEN/NOTIFY соединение с тем же search_path. При PgBouncer LISTEN требует session pooling либо прямого подключения; текущий starter использует один DSN, поэтому transaction pooling для worker не поддержан.

`/health/ready` проверяет web-процесс и БД, но не доказывает работу worker. Контролируйте worker как отдельный сервис, его exit status и рост очереди. River UI/Mailpit из Compose привязаны к loopback и предназначены для разработки. Production-панель River UI требует отдельно настроенной аутентификации и сетевого доступа.
