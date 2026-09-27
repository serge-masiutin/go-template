# Go Template

Go + Inertia + React starter для приложения с серверными маршрутами и React-интерфейсом. В одном репозитории — PostgreSQL, ORM, вход, приватные заметки, фоновые задачи, почта, AI-помощник, админка, Storybook и навыки для агента: можно начать с работающего сценария и развивать его под свой продукт.

- **Go:** `net/http`, Gonertia, GORM generics, pgx, Goose и SCS; типизированные входы и явные транзакции.
- **Фоновые задачи:** River и локальный River UI; транзакционная постановка и защита от повторных внешних эффектов.
- **Почта и AI:** go-mail + Mailpit, Genkit с Gemini/OpenAI, read-only tools, schema и лимиты; рабочие сценарии на странице Note tools.
- **Разработка:** Air, Prometheus, типизированная конфигурация; [выбор библиотек и альтернативы](docs/stack.md).
- **React:** TypeScript, Vite, Tailwind и локальный Martian Mono; навигация и формы через Inertia.
- **Доступ:** серверные сессии, CSRF, проверка владельца записей и актуальных прав администратора.
- **Проверки:** govulncheck, Go race detector, PostgreSQL integration tests, Vitest, Playwright, сборка Storybook и Docker.
- **30 skills:** 18 оригинальных Evil Martians, `clear-writing` из Rails-шаблона и 11 адаптаций под Go. Происхождение и различия сохранены в [каталоге](docs/skills.md).

```text
React Form → POST /notes → проверка CSRF и входа
           → INSERT с ID текущего пользователя
           → 303 / → Go передаёт props → React показывает заметку
```

## Начать

Нужны Git, [mise](https://mise.jdx.dev/) и запущенный Docker с Compose. Версии Go и Node закреплены в `mise.toml`; зависимости — в lock-файлах.

Создайте репозиторий кнопкой **Use this template** на GitHub, затем:

```sh
git clone git@github.com:YOUR_ACCOUNT/YOUR_APP.git
cd YOUR_APP
mise trust
mise install
mise exec -- bin/configure github.com/YOUR_ACCOUNT/YOUR_APP
mise exec -- bin/setup
mise exec -- bin/manage create-user --email you@example.com --admin
```

Последняя команда читает пароль из стандартного ввода: введите 12–72 байта и нажмите Enter. Терминал отображает ввод; для скрытого ввода используйте [команду из инструкции](docs/development.md#создать-пользователя). Пароль не передаётся аргументом процесса.

```sh
mise exec -- bin/dev
```

Откройте **http://localhost:3000**. Войдите, создайте заметку и перейдите в Admin. React обновляется через Vite HMR; Go server и worker пересобираются через Air. Страница **Note tools** отправляет письма в локальный Mailpit; для AI [задайте провайдера](docs/ai.md). Storybook запускается отдельно: `mise exec -- npm run storybook`, адрес **http://localhost:6006**.

`bin/setup` создаёт `.env` из примера, устанавливает зависимости, поднимает локальный PostgreSQL на порту 5437 и Mailpit на 1025/8025, применяет миграции и собирает frontend. Данные БД хранятся в Docker volume.

## Проверить

```sh
mise exec -- bin/ci
```

Команда проверяет Go, TypeScript, frontend-тесты, сборки приложения и Storybook, целостность skills. Для интеграционных и браузерных проверок нужна отдельная тестовая БД: [полные команды](docs/testing.md). Все эти проверки и сборка контейнера включены в [GitHub Actions](.github/workflows/ci.yml).

## Устройство

[Набор библиотек](docs/stack.md) · [Очередь и почта](docs/background.md) · [AI](docs/ai.md) · [Архитектура](docs/architecture.md) · [Разработка](docs/development.md) · [Тесты](docs/testing.md) · [Контейнер и деплой](docs/deployment.md) · [Skills и источники](docs/skills.md)

Основа — [rails-template](https://github.com/serge-masiutin/rails-template); механизмы адаптированы под Go. Realtime, загрузка файлов, восстановление пароля и публичная регистрация не установлены. Админка показывает число пользователей; управление пользователями выполняется CLI.

Код — [MIT](LICENSE). Заимствованные материалы и их лицензии перечислены в [THIRD_PARTY.md](THIRD_PARTY.md).
