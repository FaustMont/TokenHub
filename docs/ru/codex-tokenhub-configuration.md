# Подключение Codex к TokenHub: четыре метода конфигурации и восстановление

Язык: [English](../codex-tokenhub-configuration.md) | [简体中文](../zh-CN/codex-tokenhub-configuration.md) | [日本語](../ja/codex-tokenhub-configuration.md) | Русский

> Это руководство подключает локальный Codex CLI, настольное приложение Codex и IDE-расширение к TokenHub. Охватывает изолированный профиль, локальные для процесса переопределения, глобальную конфигурацию CLI и конфигурацию рабочего стола.
>
> Для кратчайшей процедуры только с профилем см. [Быстрый старт с профилем](codex-tokenhub-profile-quick-start.md).

## 1. Выбор метода

| Метод | Область | Постоянный | Лучше всего для | Восстановление |
| --- | --- | --- | --- | --- |
| Изолированный профиль | Сессии, запущенные с профилем | Да | Отдельные проекты или задачи | Запуск без `--profile tokenhub` |
| Локальное для процесса переопределение | Текущий процесс или терминал | Нет | Начальная проверка или редкое использование | Выход и очистка переменных |
| Глобальная конфигурация CLI | Локальные сессии Codex текущего пользователя | Да | CLI использует TokenHub по умолчанию | Восстановить `config.toml` |
| Конфигурация рабочего стола | Рабочий стол, CLI и IDE-расширение | Да | Долгосрочная конфигурация через приложение | Восстановить `config.toml` и перезапустить |

Используйте локальный для процесса метод для начальной проверки перед выбором постоянной конфигурации.

Codex CLI, настольное приложение и IDE-расширение используют один `~/.codex/config.toml`. Доверенные проектные файлы `.codex/config.toml` не могут переопределять `model_provider`, `model_providers` или `openai_base_url`; используйте профиль для изолированного выбора провайдера.

Токен входа в консоль TokenHub — не API-ключ проекта. Никогда не коммитьте API-ключи и не помещайте их в историю оболочки. Предпочитайте `env_key`; `experimental_bearer_token` доступен только для контролируемого использования при разработке и хранит ключ в открытом тексте.

Скриншоты должны использовать реальную конфигурацию или реальные запросы. Полностью редактируйте ключи, заголовки авторизации, токены входа и OAuth, приватные хосты, имена пользователей, пути, ID проектов и аккаунтов, ID сессий и ID запросов.

---

## 2. Предварительные условия

### 2.1 Необходимые значения

| Значение | Источник | Требование |
| --- | --- | --- |
| Base URL TokenHub | Информация о развёртывании | Заканчивается на `/v1` |
| API-ключ проекта TokenHub | **Управление ключами** в TokenHub | API-ключ проекта, не токен входа |
| ID модели | `GET /v1/models` с ключом проекта | Используйте фактический `data[].id` |

### 2.2 Задайте переменные терминала

#### macOS zsh

```bash
export TOKENHUB_BASE_URL="enter the actual TokenHub Base URL"
read -r -s "TOKENHUB_API_KEY?TokenHub project API key: "
export TOKENHUB_API_KEY
echo
```

Выполните команды по порядку. После появления приглашения `read` вставьте ключ и нажмите Enter. Символы или звёздочки не отображаются, поскольку `-s` отключает эхо ввода. `export` делает значение доступным для Codex; `echo` восстанавливает чистую строку приглашения.

#### Bash

```bash
export TOKENHUB_BASE_URL="enter the actual TokenHub Base URL"
read -r -s -p "TokenHub project API key: " TOKENHUB_API_KEY
export TOKENHUB_API_KEY
echo
```

#### Windows PowerShell

```powershell
$env:TOKENHUB_BASE_URL = Read-Host "TokenHub Base URL (must end with /v1)"
$tokenHubSecureKey = Read-Host "TokenHub project API key" -AsSecureString
$env:TOKENHUB_API_KEY = [System.Net.NetworkCredential]::new("", $tokenHubSecureKey).Password
Remove-Variable tokenHubSecureKey
```

Эти переменные действуют только в текущей сессии терминала.

### 2.3 Обнаружение моделей

#### macOS, Linux или Git Bash

```bash
curl --fail-with-body \
  --url "${TOKENHUB_BASE_URL%/}/models" \
  --header "Authorization: Bearer ${TOKENHUB_API_KEY}"
```

#### Windows PowerShell

```powershell
$tokenHubModels = Invoke-RestMethod `
  -Uri "$($env:TOKENHUB_BASE_URL.TrimEnd('/'))/models" `
  -Headers @{ Authorization = "Bearer $env:TOKENHUB_API_KEY" }

$tokenHubModels.data | Select-Object id
```

Сохраните фактический возвращённый ID модели:

```bash
read -r "TOKENHUB_MODEL_ID?Enter an actual model ID from the previous response: "
export TOKENHUB_MODEL_ID
```

PowerShell:

```powershell
$env:TOKENHUB_MODEL_ID = Read-Host "Enter an actual model ID from the previous response"
```

### 2.4 Проверка потокового Responses

Видимость модели не доказывает работоспособность маршрута. Codex требует поддержки потокового Responses API.

```bash
curl --fail-with-body --no-buffer \
  --request POST \
  --url "${TOKENHUB_BASE_URL%/}/responses" \
  --header "Authorization: Bearer ${TOKENHUB_API_KEY}" \
  --header "Content-Type: application/json" \
  --data "$(printf '{"model":"%s","input":"Reply only: Connection successful","stream":true}' "$TOKENHUB_MODEL_ID")"
```

PowerShell:

```powershell
$tokenHubRequestBody = @{
  model = $env:TOKENHUB_MODEL_ID
  input = "Reply only: Connection successful"
  stream = $true
} | ConvertTo-Json -Compress

Invoke-WebRequest `
  -Method Post `
  -Uri "$($env:TOKENHUB_BASE_URL.TrimEnd('/'))/responses" `
  -Headers @{ Authorization = "Bearer $env:TOKENHUB_API_KEY" } `
  -ContentType "application/json" `
  -Body $tokenHubRequestBody
```

Если TokenHub возвращает `provider_capability_not_supported`, администратор должен исправить маршрут модели или тип ресурса провайдера.

Для официального провайдера DeepSeek Responses и Codex ограничены моделями и доступны как для `deepseek-v4-flash`, так и для `deepseek-v4-pro`. Обе модели поддерживают серверный `web_search`, пользовательский инструмент Codex `apply_patch` и `top_logprobs` от 0 до 20, но не поддерживают изображения или файлы в качестве входных данных. Responses API DeepSeek не сохраняет состояние, поэтому клиенты должны отправлять полную историю беседы в `input` при каждом ходе вместо использования `previous_response_id` или `conversation`. DeepSeek управляет кэшированием контекста автоматически. При `TOKENHUB_CACHE_AFFINITY_ENABLED=true` TokenHub использует стабильные подсказки сессии Codex, такие как `session-id`, `client_metadata.session_id` или `prompt_cache_key`, чтобы держать последовательные ходы Responses на одном и том же вышестоящем аккаунте; этот ключ управляет маршрутизацией шлюза и не создаёт отдельный кэш ответов TokenHub.

---

## 3. Метод 1: изолированный профиль

### 3.1 Создайте резервную копию профиля

Пути к профилю:

- macOS / Linux: `~/.codex/tokenhub.config.toml`
- Windows: `%USERPROFILE%\.codex\tokenhub.config.toml`

```bash
if [ -f "$HOME/.codex/tokenhub.config.toml" ]; then
  cp -p "$HOME/.codex/tokenhub.config.toml" \
    "$HOME/.codex/tokenhub.config.toml.before-edit.$(date +%Y%m%d-%H%M%S)"
fi
```

PowerShell:

```powershell
if (Test-Path "$env:USERPROFILE\.codex\tokenhub.config.toml") {
  $tokenHubBackupTime = Get-Date -Format "yyyyMMdd-HHmmss"
  Copy-Item `
    "$env:USERPROFILE\.codex\tokenhub.config.toml" `
    "$env:USERPROFILE\.codex\tokenhub.config.toml.before-edit.$tokenHubBackupTime"
}
```

### 3.2 Запишите профиль

```toml
model_provider = "tokenhub"
model = "enter an actual model ID returned by GET /v1/models"

[model_providers.tokenhub]
name = "TokenHub"
base_url = "enter the actual TokenHub Base URL"
env_key = "TOKENHUB_API_KEY"
env_key_instructions = "Set TOKENHUB_API_KEY before starting Codex"
wire_api = "responses"
```

`base_url` должен включать `/v1`. Не дублируйте таблицу провайдера или ключи верхнего уровня. Codex 0.134.0 и выше использует отдельный файл `<profile>.config.toml`.

Для хранения ключа напрямую на контролируемой личной машине разработки удалите `env_key` и `env_key_instructions`:

```toml
[model_providers.tokenhub]
name = "TokenHub"
base_url = "enter the actual TokenHub Base URL"
experimental_bearer_token = "paste your TokenHub project API key here"
wire_api = "responses"
```

Не совмещайте `experimental_bearer_token` с `env_key`, `auth` провайдера или `requires_openai_auth`. Установите разрешения в `600` и никогда не коммитьте, не загружайте, не делитесь и не делайте скриншотов файла.

![Реальная конфигурация профиля TokenHub с редактированным Base URL](../assets/codex-profile/tokenhub-profile-config-redacted.png)

*Рисунок 1: Реальная конфигурация профиля с переменной окружения. Base URL редактирован.*

### 3.3 Запуск и проверка

```bash
codex --profile tokenhub
codex --profile tokenhub --cd "/enter/the/absolute/project/path"
codex exec --profile tokenhub --cd "/enter/the/absolute/project/path" "enter the real task"
```

Выполните `/status` и подтвердите модель, провайдер TokenHub и соответствующую успешную запись в журнале запросов TokenHub.

![Реальный статус Codex через профиль TokenHub с редактированной конфиденциальной информацией](../assets/codex-profile/codex-status-redacted.png)

*Рисунок 2: `/status` с включённым профилем; детали провайдера, заголовок окна и ID сессии редактированы.*

### 3.4 Восстановление

Используйте конфигурацию по умолчанию:

```bash
codex
```

Отключите профиль:

```bash
mv "$HOME/.codex/tokenhub.config.toml" \
  "$HOME/.codex/tokenhub.config.toml.disabled"
```

При наличии восстановите ранее существовавший профиль из его резервной копии с временной меткой.

---

## 4. Метод 2: локальное для процесса переопределение

Этот метод не изменяет ни одного файла и применяется только к текущему процессу Codex.

```bash
codex \
  -c 'model_provider="tokenhub"' \
  -c "model=\"${TOKENHUB_MODEL_ID}\"" \
  -c 'model_providers.tokenhub.name="TokenHub"' \
  -c "model_providers.tokenhub.base_url=\"${TOKENHUB_BASE_URL}\"" \
  -c 'model_providers.tokenhub.env_key="TOKENHUB_API_KEY"' \
  -c 'model_providers.tokenhub.env_key_instructions="Set TOKENHUB_API_KEY before starting Codex"' \
  -c 'model_providers.tokenhub.wire_api="responses"'
```

PowerShell:

```powershell
codex `
  -c 'model_provider="tokenhub"' `
  -c "model=`"$env:TOKENHUB_MODEL_ID`"" `
  -c 'model_providers.tokenhub.name="TokenHub"' `
  -c "model_providers.tokenhub.base_url=`"$env:TOKENHUB_BASE_URL`"" `
  -c 'model_providers.tokenhub.env_key="TOKENHUB_API_KEY"' `
  -c 'model_providers.tokenhub.env_key_instructions="Set TOKENHUB_API_KEY before starting Codex"' `
  -c 'model_providers.tokenhub.wire_api="responses"'
```

Проверьте через `/status`. Выйдите из Codex для отмены переопределения. По завершении очистите переменные:

```bash
unset TOKENHUB_BASE_URL TOKENHUB_API_KEY TOKENHUB_MODEL_ID
```

---

## 5. Метод 3: глобальная конфигурация CLI

Пути к пользовательской конфигурации:

- macOS / Linux: `~/.codex/config.toml`
- Windows: `%USERPROFILE%\.codex\config.toml`

Создайте резервную копию файла:

```bash
if [ -f "$HOME/.codex/config.toml" ]; then
  cp -p "$HOME/.codex/config.toml" \
    "$HOME/.codex/config.toml.before-tokenhub.$(date +%Y%m%d-%H%M%S)"
fi
```

Добавьте тот же блок провайдера из раздела 3.2 в `config.toml`. Изменяйте существующие ключи вместо добавления дублей. `experimental_bearer_token` разрешается только на тех же условиях, что описаны выше.

Проверьте:

```bash
codex doctor --summary
codex
```

Выполните `/status` и подтвердите запрос в журналах TokenHub.

Для восстановления восстановите резервную копию или верните `model_provider` и `model` к прежним значениям и удалите таблицу провайдера TokenHub. Полностью перезапустите Codex.

---

## 6. Метод 4: конфигурация настольного приложения Codex

Настольное приложение, CLI и IDE-расширение используют один `~/.codex/config.toml`.

Откройте:

**Настройки → Конфигурация → Открыть config.toml**

Создайте резервную копию файла и добавьте блок провайдера из раздела 3.2.

Настольные приложения, запущенные вне терминала, обычно не наследуют переменные терминала. Добавьте ключ в `~/.codex/.env`:

```dotenv
TOKENHUB_API_KEY=enter the actual TokenHub project API key
```

```bash
chmod 600 "$HOME/.codex/.env"
```

Никогда не перезаписывайте несвязанные значения `.env` и не включайте файл в систему контроля версий или скриншоты.

Полностью перезапустите Codex, создайте локальную задачу, подтвердите модель и проверьте запрос в журналах TokenHub. Локальный `config.toml` не управляет моделью по умолчанию для облачных задач Codex.

Для восстановления восстановите `config.toml`, удалите только добавленную строку `TOKENHUB_API_KEY` из `.env` и перезапустите.

---

## 7. Устранение неполадок

| Симптом | Вероятная причина | Действие |
| --- | --- | --- |
| Отсутствует `TOKENHUB_API_KEY` | Процесс не получил переменную | Проверьте переменную; перезапустите настольное приложение после обновления `.env` |
| HTTP 401 / `invalid_api_key` | Отсутствующий, некорректный или нераспознанный ключ проекта | Используйте API-ключ проекта TokenHub, а не токен входа в консоль |
| HTTP 403 | Отключённый или истёкший ключ, или запрещённая модель | Проверьте проект, ключ, белый список моделей и политику |
| HTTP 404 | Неверный Base URL или ID модели | Подтвердите `/v1` и повторно запросите `GET /v1/models` |
| HTTP 429 / `quota_exceeded` | Лимит запросов, токенов, стоимости, параллелизма или провайдера | Подождите восстановления или скорректируйте политику |
| HTTP 503 / `provider_unavailable` | Нет работоспособного маршрута | Проверьте маршрут, провайдера и работоспособность ресурса аккаунта |
| HTTP 501 / `provider_capability_not_supported` | Маршрут не может обеспечить Responses или потоковый Responses | Измените маршрут модели или ресурс провайдера |

Проверьте только существование ключа:

```bash
test -n "${TOKENHUB_API_KEY:-}" && echo "TOKENHUB_API_KEY is set"
```

Приоритет конфигурации:

1. Аргументы CLI и `--config`
2. Доверенный проектный `.codex/config.toml`
3. Выбранный файл профиля
4. Пользовательский `~/.codex/config.toml`
5. Системная конфигурация
6. Значения по умолчанию Codex

## 8. Ссылки

- [Основы конфигурации Codex](https://learn.chatgpt.com/docs/config-file/config-basic)
- [Расширенная конфигурация Codex](https://learn.chatgpt.com/docs/config-file/config-advanced)
- [Переменные окружения Codex](https://learn.chatgpt.com/docs/config-file/environment-variables)
- [Справочник конфигурации Codex](https://learn.chatgpt.com/docs/config-file/config-reference)
