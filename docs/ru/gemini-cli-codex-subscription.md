# Использование Codex Subscription GPT из Gemini CLI

TokenHub предоставляет нативный интерфейс совместимости с Gemini `v1beta`. Официальный Gemini CLI может напрямую подключаться к TokenHub, который затем маршрутизирует запрос к аккаунту OpenAI Codex Subscription. CCswitch и другие локальные прокси протоколов не требуются.

## Предварительные условия

- В TokenHub настроен работоспособный аккаунт OpenAI Codex Subscription.
- GPT-модель, например `gpt-5.5`, включена и маршрутизирована к этому провайдеру.
- Ключ проекта TokenHub разрешает эту модель.
- Установленный Gemini CLI поддерживает `GOOGLE_GEMINI_BASE_URL`.

Используйте HTTPS, кроме `localhost`, `127.0.0.1` или `[::1]`. Не добавляйте `/v1beta` к Base URL.

## Запуск без изменения существующих настроек

Следующие переменные окружения применяются только к этой команде и не редактируют `~/.gemini/settings.json`:

```bash
export TOKENHUB_GEMINI_KEY='your TokenHub project key'

GEMINI_API_KEY="$TOKENHUB_GEMINI_KEY" \
GOOGLE_GEMINI_BASE_URL='https://tokenhub.example.com' \
GEMINI_MODEL='gpt-5.5' \
gemini -m gpt-5.5

unset TOKENHUB_GEMINI_KEY
```

Для локального экземпляра TokenHub:

```bash
GEMINI_API_KEY="$TOKENHUB_GEMINI_KEY" \
GOOGLE_GEMINI_BASE_URL='http://127.0.0.1:8080' \
GEMINI_MODEL='gpt-5.5' \
gemini -m gpt-5.5
```

Gemini CLI отправляет `GEMINI_API_KEY` в заголовке `x-goog-api-key`; используйте ключ проекта TokenHub, а не OAuth-токен доступа OpenAI.

## Сохранение конфигурации для одного проекта

Создайте `.gemini/.env` внутри этого проекта:

```dotenv
GEMINI_API_KEY=your TokenHub project key
GOOGLE_GEMINI_BASE_URL=https://tokenhub.example.com
GEMINI_MODEL=gpt-5.5
```

Добавьте `.gemini/.env` в `.gitignore` этого проекта. Конфигурация Gemini на уровне пользователя при этом остаётся неизменной.

## Поддерживаемый интерфейс

- `GET /v1beta/models` и `GET /v1beta/models/{model}`
- `generateContent`, SSE `streamGenerateContent` и `countTokens`
- Текст, встроенные изображения, клиентские инструменты, результаты функций и многоходовые вызовы инструментов
- Продолжение рассуждений Codex, передаваемое в Gemini `thoughtSignature`
- Аффинность аккаунта подписки Codex в рамках одной беседы Gemini

Серверные инструменты Gemini `googleSearch`, `codeExecution` и `cachedContent` не поддерживаются. Локальные инструменты Gemini CLI — чтение файлов, команды оболочки и редактирование — остаются доступными.

Официальный справочник по `GOOGLE_GEMINI_BASE_URL`: [конфигурация Gemini CLI](https://github.com/google-gemini/gemini-cli/blob/main/docs/reference/configuration.md).
