# Фреймворк миграции TokenHub

Фреймворк миграции TokenHub предоставляет повторяемый, идемпотентный рабочий процесс для переноса конфигураций конкурирующих AI-шлюзов в TokenHub.

## Текущее состояние

Текущая ветка включает работающий канонический bundle, TokenHub sink с поддержкой выполнения как на базе локального хранилища (store), так и через удалённый Admin API, файловый адаптер для LiteLLM и готовый CLI-процесс для команд `extract`, `plan`, `apply`, `verify` и `rollback`.

## Архитектура

См. [architecture.md](./architecture.md) для получения сведений о структуре фреймворка и руководства по его расширению.

## Поддерживаемые источники

| Источник | Адаптер | Поддерживаемые версии | Статус |
|----------|---------|-----------------------|--------|
| LiteLLM | `litellm` | ≥1.52.0, <1.70.0 | Базовый |

См. [litellm.md](./litellm.md) для специфики LiteLLM.

## Канонический bundle

Промежуточное представление, используемое между адаптерами источников и TokenHub sink. См. [bundle-schema.md](./bundle-schema.md) для схемы и политики совместимости.

## CLI

```bash
tokenhub-migrate inspect litellm --from proxy_config.yaml
tokenhub-migrate extract litellm --from proxy_config.yaml --out bundle.json
tokenhub-migrate plan --bundle bundle.json
tokenhub-migrate apply --bundle bundle.json
tokenhub-migrate verify --bundle bundle.json
tokenhub-migrate rollback --checkpoint checkpoint.json
```

## Обработка секретов

Секреты в bundle хранятся в виде ссылок `{"$secretRef": "ENV_NAME"}`. Sink разрешает их во время применения (apply) из переменных окружения, файла или интерактивного ввода. В сам bundle открытые секреты не внедряются.

## Документация

- [Архитектура](./architecture.md)
- [Схема bundle](./bundle-schema.md)
- [Адаптер LiteLLM](./litellm.md)
- [Справочник по CLI](./cli.md)
- [Сквозное E2E-тестирование](./e2e.md)
