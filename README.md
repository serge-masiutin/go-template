# Go Template

Go + Inertia + React starter для приложения с серверными маршрутами и React-интерфейсом. В одном репозитории — PostgreSQL, вход, приватные заметки, админка, Storybook и навыки для агента: можно начать с работающего сценария и развивать его под свой продукт.

- **Go:** `net/http`, Gonertia, pgx и SCS; явный SQL, типизированные входы и ошибки.
- **React:** TypeScript, Vite, Tailwind и локальный Martian Mono; навигация и формы через Inertia.
- **Доступ:** серверные сессии, CSRF, проверка владельца записей и актуальных прав администратора.
- **Проверки:** Go race detector, PostgreSQL integration tests, Vitest, Playwright, сборка Storybook и Docker.
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

Откройте **http://localhost:3000**. Войдите, создайте заметку и перейдите в Admin. React обновляется через Vite HMR; после изменения Go перезапустите `bin/dev`. Storybook запускается отдельно: `mise exec -- npm run storybook`, адрес **http://localhost:6006**.

`bin/setup` создаёт `.env` из примера, устанавливает зависимости, поднимает локальный PostgreSQL на порту 5437, применяет миграции и собирает frontend. Данные БД хранятся в Docker volume.

## Проверить

```sh
mise exec -- bin/ci
```

Команда проверяет Go, TypeScript, frontend-тесты, сборки приложения и Storybook, целостность skills. Для интеграционных и браузерных проверок нужна отдельная тестовая БД: [полные команды](docs/testing.md). Все эти проверки и сборка контейнера включены в [GitHub Actions](.github/workflows/ci.yml).

## Устройство

[Архитектура](docs/architecture.md) · [Разработка](docs/development.md) · [Тесты](docs/testing.md) · [Контейнер и деплой](docs/deployment.md) · [Skills и источники](docs/skills.md)

Это основа на базе [rails-template](https://github.com/serge-masiutin/rails-template), а не перенос каждой Rails-интеграции. Почта, очередь, AI, realtime, загрузка файлов, восстановление пароля и публичная регистрация не установлены. Навыки описывают, как добавлять такие возможности при необходимости. Админка показывает число пользователей; управление пользователями выполняется CLI.

Код — [MIT](LICENSE). Заимствованные материалы и их лицензии перечислены в [THIRD_PARTY.md](THIRD_PARTY.md).
