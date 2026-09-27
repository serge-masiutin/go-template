# Заимствованные материалы

- Основа проекта и локальные исходные навыки: [serge-masiutin/rails-template](https://github.com/serge-masiutin/rails-template), MIT. Точные ревизии указаны в [skills-lock.json](config/skills-lock.json).
- Evil Martians Agent Skills: [официальный каталог](https://evilmartians.com/.well-known/agent-skills/index.json), [публичный репозиторий](https://github.com/evilmartians/agent-skills). Copyright © 2026 Evil Martians, [MIT notice](third_party/licenses/evilmartians-agent-skills.txt). Исходные React/Storybook/general skills и адаптированные Inertia skills сохраняют атрибуцию.
- `layered-go` адаптирован из `layered-rails` Владимира Дементьева (Vladimir Dementyev / palkan), через указанную ревизию Rails-шаблона. Источник: [layered-rails-skills](https://github.com/palkan/layered-rails-skills); upstream объявляет MIT в `layered-rails/.claude-plugin/plugin.json` и gemspec. Методическая основа — *Layered Design for Ruby on Rails Applications*. Go-адаптация не является официальным изданием автора или Evil Martians. Оригинальные файлы, ссылки автора и изменения сохранены в [provenance](third_party/skills).
- `clear-writing`: проектный навык из Rails-шаблона, сохранён без изменений; исходные авторские/книжные ссылки находятся внутри навыка. Это не распространяет лицензию шаблона на упомянутые книги.
- Martian Mono: [Evil Martians](https://github.com/evilmartians/mono), SIL Open Font License 1.1; [полный текст](web/public/fonts/OFL.txt).

Зависимости Go и npm имеют собственные лицензии. Манифесты и lock-файлы определяют точные используемые версии. Архивы в `third_party/skills` сохраняют исходные инструкции для ревью и сравнения; активные навыки находятся только в `.agents/skills`.
