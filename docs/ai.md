# AI-помощник

Genkit Go выполняет один типизированный flow `notes_assistant_v1` с read-only tool `list_notes`. HTTP создаёт операцию и River job; worker получает ID, загружает владельца и вопрос, вызывает flow, затем сохраняет ответ. Страница `/tools` показывает пять последних запросов и опрашивает их состояние, пока работа не завершена.

## Включение

Приложение запускается с `AI_ENABLED=false` без API key. Для включения задайте в локальном `.env` или secrets deployment:

```sh
AI_ENABLED=true
AI_MODEL=googleai/gemini-3.8-flash
AI_API_KEY=your-provider-key
AI_TIMEOUT=30s
AI_MAX_TURNS=3
AI_MAX_OUTPUT_TOKENS=1024
```

Указанный ID используется в локальных protocol tests; доступность конкретной модели в вашем аккаунте проверяется у провайдера. Для OpenAI используйте префикс `openai/` и ID Chat Completions модели с tools и structured output. Genkit-плагин в этой версии работает через Chat Completions; Responses-only модели ему не подходят. Ключ относится к выбранному провайдеру. После изменения конфигурации перезапустите server и worker.

При включении обязательны модель и ключ. Ошибка конфигурации останавливает процесс до HTTP-запроса к провайдеру. Создание клиента и запуск приложения не запускают генерацию. Пользователь видит перед отправкой, что его вопрос и заметки будут переданы AI-провайдеру.

## Контракт и границы

- Вопрос: 1–500 Unicode-символов, JSON ограничен HTTP middleware.
- Tool не принимает owner ID. На каждый вызов flow создаётся `ai.NewTool`, замыкающий доверенный ID операции. Чтение ограничено 50 последними заметками владельца, по 2000 символов каждая.
- У tool нет записи, SMTP, shell, browser или произвольных URL. Текст заметок — данные, а не инструкции.
- Prompt: `internal/assistant/prompts/notes-v1.txt`; `PromptVersion` хранится вместе с моделью в операции.
- Output: обязательный `answer`, непустой plain text до 4000 символов. Успешный ответ без вызова инструмента отклоняется. React выводит текст, не HTML модели.
- `AI_TIMEOUT`: 1s–5m, default 30s на весь flow; `AI_MAX_TURNS`: 1–10; `AI_MAX_OUTPUT_TOKENS`: 128–8192 на модельный вызов. Общая стоимость зависит от числа turns и входных заметок.
- OpenAI SDK retries выключены. У используемого Gemini generate path нет автоматических HTTP retries; это проверяется тестом 503. River mail/AI также имеют одну попытку.
- Raw provider errors не сохраняются в River и логах. В приложении остаётся обёрнутая причина для диагностики безопасными категориями. Тексты вопросов/ответов находятся только в owner-scoped таблице, не в queue payload или метриках.

Не подключайте Genkit telemetry/reflection к production-процессу по примеру dev quickstart: такие инструменты могут сохранять полный контент. `GENKIT_ENV`, `GENKIT_TELEMETRY_SERVER`, `GENKIT_REFLECTION_V2_SERVER` при включённом AI отклоняются. OTLP export, conversation memory, embeddings и multi-agent orchestration не настроены. Срок хранения вопросов/ответов определяется вашей продуктовой политикой; starter не удаляет их автоматически, кроме удаления аккаунта.

## Проверки

`go test ./internal/assistant` проверяет schema, отсутствие доказательств, loop limit, provider outage, границу длины и HTTP-контракты обоих SDK, включая Gemini thought signature. PostgreSQL integration tests проверяют атомарность enqueue, изоляцию заметок, отмену после удаления аккаунта и отсутствие повторной генерации.

Эти проверки работают на синтетических данных и не тратят provider credits. Они доказывают контракты приложения, а не качество рассуждений реальной модели. Перед публикацией AI-функции выполните [набор оценок](../evals/notes-assistant.md) на выбранной модели; зафиксируйте её ID, prompt version, результаты, latency и стоимость по данным провайдера. Отсутствующий usage не считайте нулём.
