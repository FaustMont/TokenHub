# Подключение Codex к TokenHub: быстрый старт с профилем

Язык: [English](../codex-tokenhub-profile-quick-start.md) | [简体中文](../zh-CN/codex-tokenhub-profile-quick-start.md) | [日本語](../ja/codex-tokenhub-profile-quick-start.md) | Русский

> Это руководство для пользователей, которым нужно только подключить Codex к TokenHub через изолированный профиль `tokenhub`. Охватывает создание профиля, настройку API-ключа, проверку и восстановление.
>
> Для сравнения профиля, локального для процесса, глобального CLI и конфигурации рабочего стола см. [Подключение Codex к TokenHub: четыре метода конфигурации и восстановление](codex-tokenhub-configuration.md).

Прежде чем пользователи будут следовать этому руководству, администратор должен настроить провайдера OpenAI Codex, ресурсы аккаунта подписки и маршруты модели в TokenHub, а затем создать API-ключ для проекта.

## 1. Как работает профиль

Профиль `tokenhub` загружается только когда Codex запущен с `--profile tokenhub`.

| Команда | Основная конфигурация | Путь запроса |
| --- | --- | --- |
| `codex` | `~/.codex/config.toml` | Провайдер Codex по умолчанию |
| `codex --profile tokenhub` | `~/.codex/tokenhub.config.toml` | TokenHub |

Профиль не перезаписывает `config.toml` по умолчанию, влияет только на явно выбранные сессии и может быть обойдён пропуском аргумента профиля.

Скриншоты должны быть получены из реальной конфигурации или реальных запросов. Полностью редактируйте API-ключи, токены входа, OAuth-токены, идентификаторы аккаунтов, ID сессий, приватные хосты, имена пользователей, абсолютные пути, имена проектов и ID запросов.

---

## 2. Использование существующего профиля

Перед началом убедитесь, что TokenHub доступен, API-ключ проекта действителен, а настроенная модель имеет работоспособный маршрут.

Для профиля, использующего `env_key`, введите ключ без эха в терминале и запустите Codex:

```bash
read -r -s "TOKENHUB_API_KEY?TokenHub project API key: "
export TOKENHUB_API_KEY
echo

codex --profile tokenhub
```

Если профиль использует `experimental_bearer_token`, пропустите `read` и `export` и запустите `codex --profile tokenhub` напрямую.

После запуска выполните `/status`. `Model provider` должен показывать TokenHub, а `Model` — совпадать с профилем.

---

## 3. Создание профиля

### 3.1 Проверьте Codex CLI

```bash
codex --version
```

Если команда недоступна, сначала установите и войдите в Codex CLI.

### 3.2 Проверьте TokenHub

В этом руководстве используется локальный эндпоинт:

```text
http://127.0.0.1:8080
```

Выполните проверку работоспособности:

```bash
curl --fail-with-body http://127.0.0.1:8080/healthz
```

Ожидаемый ответ:

```json
{"service":"tokenhub-backend","status":"ok"}
```

Если сервис недоступен, запустите TokenHub из репозитория:

```bash
cd "/enter/the/absolute/path/to/TokenHub"
./start.sh
```

Оставьте процесс запущенным и продолжите в новом терминале.

### 3.3 Создайте резервную копию существующего профиля

```bash
mkdir -p "$HOME/.codex"

if [ -f "$HOME/.codex/tokenhub.config.toml" ]; then
  cp -p "$HOME/.codex/tokenhub.config.toml" \
    "$HOME/.codex/tokenhub.config.toml.before-edit.$(date +%Y%m%d-%H%M%S)"
fi
```

### 3.4 Запишите профиль

Откройте файл:

```bash
nano "$HOME/.codex/tokenhub.config.toml"
```

Рекомендуемая конфигурация с переменной окружения:

```toml
model_provider = "tokenhub"
model = "gpt-5.6-luna"

[model_providers.tokenhub]
name = "TokenHub Local"
base_url = "http://127.0.0.1:8080/v1"
env_key = "TOKENHUB_API_KEY"
env_key_instructions = "Set TOKENHUB_API_KEY before starting Codex"
wire_api = "responses"
```

В `nano` нажмите `Control-O`, подтвердите имя файла нажатием Enter и нажмите `Control-X`.

Требования:

- `base_url` должен включать `/v1`.
- Не объявляйте `model_provider`, `model` или `[model_providers.tokenhub]` более одного раза.
- Codex 0.134.0 и выше использует отдельный файл `<profile>.config.toml` вместо устаревшей таблицы `[profiles.tokenhub]`.

Для контролируемой личной машины разработки вы можете хранить API-ключ проекта напрямую. Удалите `env_key` и `env_key_instructions` и используйте:

```toml
[model_providers.tokenhub]
name = "TokenHub Local"
base_url = "http://127.0.0.1:8080/v1"
experimental_bearer_token = "paste your TokenHub project API key here"
wire_api = "responses"
```

Не совмещайте `experimental_bearer_token` с `env_key`, `[model_providers.tokenhub.auth]` или `requires_openai_auth`. Ключ хранится в открытом тексте, и этот вариант предназначен только для использования при разработке.

```bash
chmod 600 "$HOME/.codex/tokenhub.config.toml"
```

Никогда не коммитьте, не загружайте, не делитесь и не делайте скриншотов профиля, содержащего ключ.

![Реальная конфигурация профиля TokenHub с редактированным Base URL](../assets/codex-profile/tokenhub-profile-config-redacted.png)

*Рисунок 1: Реальная конфигурация профиля с переменной окружения. Base URL редактирован, API-ключ в файле не хранится.*

### 3.5 Необязательно: проверьте TOML

Пропустите этот шаг, если профиль уже успешно запускается и проходит тест подключения в разделе 4.3.

```bash
python3 - <<'PY'
from pathlib import Path
import tomllib

path = Path.home() / ".codex" / "tokenhub.config.toml"
with path.open("rb") as file:
    config = tomllib.load(file)

print("Profile configuration loaded")
print("Model:", config["model"])
print("Provider:", config["model_provider"])
print("Base URL:", config["model_providers"]["tokenhub"]["base_url"])
PY
```

Ожидаемый вывод для этой среды:

```text
Profile configuration loaded
Model: gpt-5.6-luna
Provider: tokenhub
Base URL: http://127.0.0.1:8080/v1
```

---

## 4. Ежедневное использование и проверка

### 4.1 Введите API-ключ проекта

Этот раздел применяется только когда профиль использует `env_key = "TOKENHUB_API_KEY"`.

1. Выполните команду `read` ниже.
2. Когда появится приглашение, вставьте реальный API-ключ проекта и нажмите Enter.
3. Терминал не отображает символы или звёздочки, поскольку эхо ввода отключено.
4. Выполните `export TOKENHUB_API_KEY`, затем `echo`.

```bash
read -r -s "TOKENHUB_API_KEY?TokenHub project API key: "
export TOKENHUB_API_KEY
echo
```

- `read` считывает одну строку из терминала.
- `-r` сохраняет обратные косые черты.
- `-s` отключает эхо ввода.
- `TOKENHUB_API_KEY?...` сохраняет значение в `TOKENHUB_API_KEY` и отображает текст после `?`.
- `export TOKENHUB_API_KEY` делает значение доступным для Codex.
- `echo` восстанавливает чистую строку приглашения.

Переменная существует только в текущей сессии терминала. Избегайте помещать реальный ключ напрямую в команду `export`, поскольку он может попасть в историю оболочки.

Проверьте только существование переменной:

```bash
test -n "${TOKENHUB_API_KEY:-}" &&
  echo "TOKENHUB_API_KEY is set"
```

### 4.2 Запустите Codex

В текущей директории:

```bash
codex --profile tokenhub
```

Для конкретного проекта:

```bash
codex --profile tokenhub \
  --cd "/enter/the/absolute/project/path"
```

### 4.3 Выполните одноразовый тест подключения

```bash
codex exec \
  --profile tokenhub \
  --ephemeral \
  --sandbox read-only \
  "Do not use tools. Reply only: Connection successful"
```

Реальный тест для этой среды вернул:

```text
OpenAI Codex v0.145.0
model: gpt-5.6-luna
provider: tokenhub
Connection successful
```

Успех требует ожидаемой модели, `provider: tokenhub`, финального ответа и соответствующей записи HTTP 200 в журналах запросов TokenHub.

### 4.4 Проверьте статус в реальном времени

Выполните `/status` в Codex.

![Реальный статус Codex через профиль TokenHub с редактированной конфиденциальной информацией](../assets/codex-profile/codex-status-redacted.png)

*Рисунок 2: Реальный вывод `/status`. Детали провайдера, заголовок окна и ID сессии редактированы.*

---

## 5. Путь запроса и зависимости

```text
Текущий терминал
  → профиль tokenhub
  → http://127.0.0.1:8080/v1
  → аутентификация проекта TokenHub и маршрутизация модели
  → подключённый ресурс аккаунта OpenAI Codex
  → ответ модели
```

`GET /v1/models` доказывает только что модель видима для ключа. Успешный запрос также требует работоспособного провайдера, работоспособного ресурса аккаунта, включённого маршрута модели, действительного API-ключа проекта и поддержки потокового Responses API.

В этой среде выполнен реальный потоковый запрос Responses для `gpt-5.6-luna`; TokenHub зафиксировал HTTP 200.

---

## 6. Восстановление или отключение профиля

Используйте конфигурацию Codex по умолчанию:

```bash
codex
```

Очистите ключ из текущего терминала:

```bash
unset TOKENHUB_API_KEY
```

Отключите профиль:

```bash
mv "$HOME/.codex/tokenhub.config.toml" \
  "$HOME/.codex/tokenhub.config.toml.disabled"
```

Повторно включите его:

```bash
mv "$HOME/.codex/tokenhub.config.toml.disabled" \
  "$HOME/.codex/tokenhub.config.toml"
```

Если профиль существовал до этой настройки, восстановите его резервную копию вместо замены на переименованный файл.

---

## 7. Устранение неполадок

| Симптом | Вероятная причина | Действие |
| --- | --- | --- |
| Отсутствует `TOKENHUB_API_KEY` | Текущий терминал не экспортировал ключ | Повторите раздел 4.1 |
| HTTP 401 | Отсутствующий, истёкший или неверный API-ключ проекта | Скопируйте или ротируйте ключ проекта и введите его снова |
| HTTP 503 / `provider_unavailable` | Нет работоспособного маршрута для модели | Проверьте провайдера, ресурс аккаунта, маршрут и статус провайдера |
| Профиль не найден | Файл отсутствует или неправильно назван | Проверьте `~/.codex/tokenhub.config.toml` |
| Старый провайдер всё ещё активен | Текущий процесс не перезагрузил конфигурацию | Полностью выйдите из Codex и перезапустите с профилем |
| `doctor` отклоняет `--profile` | Эта команда Codex CLI не принимает аргумент профиля | Проверьте реальным запросом `codex exec --profile tokenhub` |

Имя файла должно быть `tokenhub.config.toml`, а не `tokenhub.toml` или `tokenhub.config.toml.txt`.

---

## 8. Контроль безопасности

- Предпочитайте `env_key` и переменную окружения.
- Используйте `experimental_bearer_token` только на контролируемой личной машине разработки.
- Установите разрешения профиля в `600`, если он содержит ключ.
- Никогда не коммитьте `.env`, API-ключи, OAuth-токены, учётные данные аккаунта или профиль с bearer-токеном.
- Немедленно ротируйте ключ, если он появился в чате, скриншотах или истории оболочки.

## 9. Связанная документация

- [Подключение Codex к TokenHub: четыре метода конфигурации и восстановление](codex-tokenhub-configuration.md)
- [Руководство пользователя по API моделей](user-guide.md)
