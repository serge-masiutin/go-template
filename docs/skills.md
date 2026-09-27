# Skills: состав и происхождение

В `.agents/skills` находится 31 навык. Они входят в GitHub template и доступны агенту вместе с проектом. Отдельного сайта, discovery endpoint и установщика, меняющего глобальные настройки агента, нет.

## Исходные Evil Martians

18 навыков взяты из официального каталога EM и сохранены побайтово: `design-system`, `good-readme`, `intent-log`, `llms-visibility`, `secure-npm-package`, `skills-visibility`, `tailwind-best-practices` и 11 `sb-*` — audit, explore, figma, flows, health, hub, inventory, setup, ship, stories, wrappers.

Оригинальные React/Storybook workflows восстановлены вместо проектных Hotwire/Lookbook-переделок. URL, SHA-256 загруженных материалов и каждого файла — в [skills-lock.json](../config/skills-lock.json). Проверка в CI не позволяет незаметно менять их под Go. Ограничения проекта хранятся в корневом `AGENTS.md`.

`clear-writing` сохранён побайтово из `rails-template` на commit, указанном в lock-файле. Один `SKILL.md` выбирает полную русскую или английскую версию по языку целевого текста. В каждой версии есть руководство, шесть глав, приёмы, шпаргалка, словарь и контрольные случаи; язык запроса сам по себе не переводит редактируемый текст. При обновлении копируйте всю папку и вместе обновляйте source commit, архив, mapping и hashes в lock-файле.

## Book-to-skill

[Оригинальный skill](../.agents/skills/book-to-skill/SKILL.md) преобразует книги и документы в навыки агента. Полный upstream-комплект сохранён без изменений: инструкции, Python-экстрактор, tools, документация, тесты и [MIT-лицензия](../.agents/skills/book-to-skill/LICENSE.md). Источник — [virgiliojr94/book-to-skill, commit 80ae087](https://github.com/virgiliojr94/book-to-skill/tree/80ae087784ddbc21dbbfde355fe5509631e0e322); архив и SHA-256 всех файлов закреплены в lock-файле.

В Codex вызовите `$book-to-skill` и укажите путь к документу. Экстрактору нужен Python 3.9+. Проверка доступных обработчиков из корня проекта:

```sh
mise exec -- python3 .agents/skills/book-to-skill/scripts/extract.py --check
```

Python-пакеты для отдельных форматов устанавливаются по необходимости и не входят в runtime приложения. Для MOBI/AZW нужен Calibre; режим технических PDF использует Docling. Подробности — в [инструкции источника](../.agents/skills/book-to-skill/docs/install.md).

## Адаптации

| Исходный навык | Навык Go | Сохранённая основа и содержательная замена |
| --- | --- | --- |
| `layered-rails` | `layered-go` | Все 57 исходных документов сопоставлены с Go-версиями; дополнительно добавлены источники и контракты установленных библиотек. Сохраняются слои, specification test, критерии выделения объектов, anti-patterns, workflows и примеры рефакторинга. |
| `inertia-rails-architecture` | `inertia-go-architecture` | Серверная навигация, владение состоянием, decision trees; transport и persistence заменены на Go. |
| `inertia-rails-controllers` | `inertia-go-controllers` | Render, shared props, authorization, loading; реальные сигнатуры Gonertia. |
| `inertia-rails-forms` | `inertia-go-forms` | React Form/useForm, uploads и многошаговые формы; строгий JSON, CSRF, 303, ограничения error bags и Precognition. |
| `inertia-rails-pages` | `inertia-go-pages` | React layouts/navigation/deferred/scroll; явный resolver и контракты Gonertia. |
| `inertia-rails-setup` | `inertia-go-setup` | Проверка существующего стека; Go module, Vite, TypeScript и явная конфигурация. |
| `inertia-rails-testing` | `inertia-go-testing` | HTTP, формы, доступ, partial/deferred; `httptest`, PostgreSQL и Playwright. |
| `inertia-rails-typescript` | `inertia-go-typescript` | React-типы и Inertia augmentation; ручные DTO вместо отсутствующего Ruby-генератора. |
| `shadcn-inertia` | `shadcn-inertia` | React-интеграция компонентов; Go CSRF, типы и совместимость с CSP. Установка shadcn остаётся отдельным решением. |
| `rails-serialization` | `go-serialization` | Entity/page/shared/loading contracts; Go DTO и явная сериализация. |
| `rails-boot-profiling` | `go-boot-profiling` | Baseline, сужение области, внешние операции, deep profile, экспорт; inittrace, pprof и измерение readiness вместо require-profiler. |

Inertia/React-основа взята из сохранённых оригиналов EM в истории `rails-template`, до Hotwire-адаптации. Точные commit/path каждого источника закреплены в lock-файле. Из forms/pages удалены только отдельные Vue/Svelte-примеры: starter выбран для React. Их исходники остаются в архиве provenance.

В [third_party/skills](../third_party/skills) для каждой адаптации лежат исходный архив и полный patch. `mapping` в lock-файле сопоставляет исходные и конечные файлы, включая десять переименованных Rails-механизмов в `layered-go`. Это позволяет проверить, что изменилось, и продолжать адаптацию от конкретного оригинала.

## Go-источники и границы рекомендаций

Материалы проверены 27 сентября 2026 года, runtime — Go 1.27.1. [Список первичных источников](../.agents/skills/layered-go/references/go-sources.md) включает Go team, Sameer Ajmani, Damien Neil, Jonathan Amsterdam, Vlad Saioc, Alex Edwards, Dave Cheney и Ardan Labs. Дата публикации отделена от даты проверки: полезные старые принципы не выдаются за новости 2026 года.

`layered-go` учитывает установленный GORM, но не требует universal repository, интерфейса для каждого типа или папки для каждого концептуального слоя. Callbacks заменены явными операциями и транзакциями; `Current` — явными actor/tenant-параметрами; relations/scopes — параметризованными запросами; concerns — композицией. Для ошибок и конкурентности учтены фактические контракты `context`, pgx и текущего Go.

Skills содержат примеры расширений, а не перечень установленных функций. Почта, AI и workers реализованы через go-mail, Genkit и River; актуальные границы описаны в [наборе библиотек](stack.md). Примеры других каналов, универсального outbox, realtime, uploads и shadcn требуют отдельной реализации под задачу. Новые возможности frontend-библиотеки не означают автоматической поддержки серверным адаптером.

## Обновление

Сначала меняйте нужный навык и проверяйте поведение на реалистичном сценарии. Для локальных адаптаций после ревью выполните `python3 scripts/update-skill-lock.py` (только стандартная библиотека Python 3), затем `mise exec -- bin/ci`. Скрипт сохраняет исходные архивы, пересчитывает hashes и patch; mapping удалённых/переименованных файлов меняется явно.

Для upstream-навыка получите оригинал из официального каталога, проверьте digest, изучите изменения и обновите URL/digest/file hashes вместе. Не используйте пересчёт hashes как замену ревью. CI проверяет целостность, ссылки и fences; семантическую корректность инструкций нужно дополнительно проверять по API и рабочим сценариям.
