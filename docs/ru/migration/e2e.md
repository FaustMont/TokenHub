# Сквозные (E2E) тесты миграции

## Обзор

E2E-тестовый стенд запускает реальные сервисы LiteLLM и TokenHub и проверяет цикл миграции через удалённый Admin API против мокового OpenAI-совместимого апстрима. Mailpit обеспечивает локальную SMTP-доставку для процесса импорта пользователей со сбросом пароля.

## Предварительные требования

- Docker и Docker Compose
- Node.js 20+
- Go toolchain, соответствующий репозиторию

## Локальный запуск

```bash
cd backend && go build -o tokenhub-migrate ./cmd/tokenhub-migrate/
docker compose -f deploy/docker-compose.migration-e2e.yml config
docker compose -f deploy/docker-compose.migration-e2e.yml up -d --wait
cd sdk/migration-e2e
npm ci
TOKENHUB_MIGRATE_BIN=../../backend/tokenhub-migrate npm run test:litellm
docker compose -f ../../deploy/docker-compose.migration-e2e.yml down -v
```

## Ресурсы фикстур

- Фикстура Compose: `deploy/litellm-config.yaml`
- Конфигурация мокового апстрима: `deploy/mock-upstream.conf`
- Фикстура извлечения: `sdk/migration-e2e/fixtures/proxy_config.yaml`
- Стенд запуска: `sdk/migration-e2e/litellm-e2e.mjs`

## Текущий охват

В настоящее время стенд подтверждает:
1. Стек LiteLLM запускается из зафиксированной фикстуры
2. LiteLLM может ответить на моковый запрос chat-completion
3. TokenHub принимает извлечённый bundle через реальный Admin API
4. `verify` проходит успешно после apply, а повторный apply создаёт и обновляет ровно 0 ресурсов
5. Файлы контрольной точки (checkpoint) и одноразовых API-ключей сохраняются
6. Rollback удаляет созданные ресурсы, включая импортированного пользователя, а последующий verify обнаруживает откаченное состояние

## CI

Рабочий процесс выполняет модульные проверки миграции при соответствующих изменениях в бэкенде, документации, SDK, развёртывании и workflow. E2E-задание запускается при push и на PR с меткой `migration:e2e`.

## Устранение неполадок

- Убедитесь, что демон Docker запущен
- Проверьте доступность портов 4000, 8080 и 8081
- Проверьте логи сервисов через `docker compose logs` при ошибках
