# CLI tokenhub-migrate

## Команды

| Команда | Описание |
|---------|----------|
| `sources` | Список зарегистрированных адаптеров источников |
| `inspect [source]` | Проверить конфигурацию исходного шлюза |
| `extract [source]` | Извлечь канонический bundle миграции |
| `plan` | Пробный запуск: показать, что сделает apply на удалённом экземпляре TokenHub |
| `apply` | Применить bundle к удалённому экземпляру TokenHub через Admin API |
| `verify` | Проверить согласованность bundle на удалённом экземпляре TokenHub |
| `rollback` | Выполнить откат из файла контрольной точки (checkpoint) на удалённом экземпляре TokenHub |

Команды `plan`, `apply`, `verify` и `rollback` требуют указания целевого экземпляра TokenHub: передайте `--to` или установите `TOKENHUB_API` (а также `--token` или `TOKENHUB_ADMIN_TOKEN`). CLI отказывается выполнять эти команды против временного хранилища в памяти и завершается с кодом 5.

## Общие флаги

| Флаг | Описание | По умолчанию |
|------|----------|--------------|
| `--secret-source` | Источник разрешения секретов: env, file | `env` |
| `--secret-file` | Файл формата `key=value`, используемый с `--secret-source=file` | — |
| `--id-strategy` | Стратегия генерации ID: stable, prefixed, source | `prefixed` |
| `--to` | Базовый URL Admin API TokenHub (или `TOKENHUB_API`) | — |
| `--token` | Токен Admin API (или `TOKENHUB_ADMIN_TOKEN`) | — |
| `--report` | Зарезервировано для структурированного вывода отчётов | — |
| `--log-level` | Зарезервировано для управления уровнем логирования | `info` |

### Выходные файлы `apply`

| Флаг | Описание | По умолчанию |
|------|----------|--------------|
| `--checkpoint-out` | JSON контрольной точки отката (записывается с правами 0600) | `<bundle>.checkpoint.json` |
| `--new-keys-out` | JSON с секретами вновь созданных API-ключей (права 0600, открытый текст виден один раз — безопасно распространите, затем удалите) | `<bundle>.new-keys.json` |

> Примечание: `apply` создаёт пользователей через эндпоинт CSV-импорта Admin API, который требует наличия активного канала email-уведомлений на целевом экземпляре. Применение bundle, создающего новых пользователей, завершится ошибкой без него, и каждый вновь импортированный пользователь получает email со сбросом пароля во время apply.
>
> Удалённое применение (apply) не является транзакционным. Если последующий ресурс завершается ошибкой после изменения предыдущих ресурсов, команда всё равно сохраняет частичную контрольную точку отката и все одноразовые API-ключи перед возвратом кода завершения 5.

## Коды завершения

| Код | Значение |
|-----|----------|
| 0 | Успех |
| 3 | Несоответствие при verify |
| 4 | Источник недоступен для чтения |
| 5 | Sink отклонён |
| 6 | Несоответствие схемы bundle |

## Пошаговое руководство LiteLLM

```bash
# Проверить конфигурацию LiteLLM
tokenhub-migrate inspect litellm --from proxy_config.yaml

# Извлечь bundle
tokenhub-migrate extract litellm --from proxy_config.yaml --out bundle.json

# Опционально: сохранить исходные ID вместо используемых по умолчанию prefixed ID
tokenhub-migrate extract litellm --from proxy_config.yaml --out bundle.json --id-strategy source

# Указать целевой экземпляр TokenHub для последующих команд
export TOKENHUB_API=http://localhost:8080
export TOKENHUB_ADMIN_TOKEN=<admin-token>

# Спланировать миграцию
tokenhub-migrate plan --bundle bundle.json

# Выполнить apply (пробный запуск)
tokenhub-migrate apply --bundle bundle.json --dry-run

# Выполнить реальное применение; сохраняет bundle.json.checkpoint.json и, если
# генерировались API-ключи, bundle.json.new-keys.json (оба с правами 0600)
tokenhub-migrate apply --bundle bundle.json

# Разрешить ссылки на секреты bundle из файла key=value
tokenhub-migrate apply --bundle bundle.json --secret-source file --secret-file migration.secrets

# Проверить поведение команды
tokenhub-migrate verify --bundle bundle.json

# Выполнить откат при необходимости
tokenhub-migrate rollback --checkpoint bundle.json.checkpoint.json
```
