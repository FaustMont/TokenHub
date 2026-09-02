# Руководство по настройке PostgreSQL

TokenHub поддерживает PostgreSQL в качестве производственной базы данных. В этом руководстве объясняется, как настроить и развернуть PostgreSQL.

## Почему PostgreSQL?

- **Производственные среды** — PostgreSQL — это корпоративная реляционная СУБД, подходящая для высококонкурентных сценариев
- **Целостность данных** — более надёжная поддержка транзакций и управление параллелизмом
- **Масштабируемость** — поддерживает горизонтальное масштабирование и репликацию основной базы
- **Резервное копирование и восстановление** — зрелая экосистема инструментов для резервного копирования

SQLite остаётся вариантом по умолчанию, подходящим для:
- Разработки и тестовых сред
- Небольших развёртываний (<1000 пользователей)
- Простых требований к развёртыванию

## Быстрый старт

### Использование Docker Compose

1. **Скопируйте конфигурацию переменных окружения**

```bash
cp deploy/.env.example deploy/.env
```

2. **Отредактируйте файл .env и задайте пароль PostgreSQL**

```bash
POSTGRES_PASSWORD=your-secure-password
TOKENHUB_SECRET_KEY=your-secret-key
# Необязательно: TOKENHUB_ADMIN_TOKEN=your-admin-token
# Необязательно: TOKENHUB_BOOTSTRAP_ADMIN_PASSWORD=your-initial-password
```

3. **Запустите сервисы**

```bash
docker compose --env-file deploy/.env -f deploy/docker-compose.postgres.yml up -d
```

4. **Откройте приложение**

- Фронтенд: http://localhost:3000
- Backend API: http://localhost:8080
- Проверка работоспособности: http://localhost:8080/healthz

Аккаунт администратора по умолчанию:
- Логин: `admin`
- Пароль: настроенный `TOKENHUB_BOOTSTRAP_ADMIN_PASSWORD`. Если он не задан, получите сгенерированное значение командой `tokenhub initial-admin-password` внутри бэкенд-контейнера.

Смените начальный пароль после входа. Развёртывания PostgreSQL должны всегда предоставлять один и тот же стабильный `TOKENHUB_SECRET_KEY` каждой реплике; TokenHub не генерирует этот общий ключ.

### Ручная установка PostgreSQL

1. **Установите PostgreSQL**

macOS:
```bash
brew install postgresql@16
brew services start postgresql@16
```

Ubuntu/Debian:
```bash
sudo apt install postgresql-16
sudo systemctl start postgresql
```

2. **Создайте базу данных и пользователя**

```bash
sudo -u postgres psql
```

```sql
CREATE USER tokenhub WITH PASSWORD 'your-password';
-- Сделайте tokenhub владельцем базы данных, чтобы он мог создавать таблицы в публичной схеме.
-- В PostgreSQL 15/16 GRANT ALL PRIVILEGES ON DATABASE одного НЕ ДОСТАТОЧНО —
-- он не предоставляет CREATE на публичной схеме, что приводит к ошибке
-- GORM AutoMigrate «permission denied for schema public».
CREATE DATABASE tokenhub OWNER tokenhub;
GRANT ALL PRIVILEGES ON DATABASE tokenhub TO tokenhub;
\q
```

Если база данных уже существует и принадлежит другой роли (например, `postgres`),
подключитесь к ней и предоставьте привилегии схемы явно:

```sql
\c tokenhub
GRANT ALL ON SCHEMA public TO tokenhub;
ALTER SCHEMA public OWNER TO tokenhub;
\q
```

3. **Настройте переменные окружения**

Задайте следующее в `backend/.env`:

```bash
TOKENHUB_DATABASE_URL=postgresql://tokenhub:your-password@localhost:5432/tokenhub?sslmode=disable
TOKENHUB_DB_MAX_OPEN_CONNS=25
TOKENHUB_DB_MAX_IDLE_CONNS=5
TOKENHUB_DB_CONN_MAX_LIFETIME_MINUTES=30
```

4. **Запустите бэкенд**

```bash
cd backend
go run ./cmd/tokenhub
```

## Конфигурация пула соединений

PostgreSQL поддерживает конфигурацию пула соединений; настройте его в соответствии с нагрузкой:

| Переменная окружения | По умолчанию | Описание |
|---------|--------|------|
| `TOKENHUB_DB_MAX_OPEN_CONNS` | 25 | Максимальное количество открытых соединений |
| `TOKENHUB_DB_MAX_IDLE_CONNS` | 5 | Максимальное количество незанятых соединений |
| `TOKENHUB_DB_CONN_MAX_LIFETIME_MINUTES` | 30 | Максимальное время жизни соединения (минуты) |

**Рекомендуемые конфигурации:**

- **Малый масштаб (<100 пользователей)**: MaxOpenConns=10, MaxIdleConns=2
- **Средний масштаб (100–1000 пользователей)**: MaxOpenConns=25, MaxIdleConns=5 (по умолчанию)
- **Крупный масштаб (>1000 пользователей)**: MaxOpenConns=50, MaxIdleConns=10

## Формат строки подключения к базе данных

```
postgresql://[user[:password]@][host][:port][/dbname][?param1=value1&...]
```

Примеры:

```bash
# Локальная разработка
postgresql://tokenhub:password@localhost:5432/tokenhub?sslmode=disable

# Производство (SSL включён)
postgresql://tokenhub:password@db.example.com:5432/tokenhub?sslmode=require

# Параметры пула соединений
postgresql://user:pass@host:5432/db?pool_max_conns=25&pool_min_conns=5
```

## Резервное копирование и восстановление

Встроенная функция резервного копирования TokenHub поддерживает только SQLite. Для PostgreSQL используйте `pg_dump` и `pg_restore`.

### Резервное копирование базы данных

```bash
pg_dump -h localhost -U tokenhub -d tokenhub -F c -f tokenhub_backup_$(date +%Y%m%d).dump
```

### Восстановление базы данных

```bash
pg_restore -h localhost -U tokenhub -d tokenhub -c tokenhub_backup_20260721.dump
```

### Автоматизированное резервное копирование (Cron)

Создайте скрипт резервного копирования `/usr/local/bin/backup-tokenhub.sh`:

```bash
#!/bin/bash
BACKUP_DIR="/var/backups/tokenhub"
DATE=$(date +%Y%m%d_%H%M%S)
mkdir -p $BACKUP_DIR

pg_dump -h localhost -U tokenhub -d tokenhub -F c -f $BACKUP_DIR/tokenhub_$DATE.dump

# Хранить резервные копии за последние 7 дней
find $BACKUP_DIR -name "tokenhub_*.dump" -mtime +7 -delete
```

Добавьте в crontab (резервное копирование ежедневно в 2:00):

```bash
0 2 * * * /usr/local/bin/backup-tokenhub.sh
```

## Миграция с SQLite на PostgreSQL

Текущая версия TokenHub не включает автоматический инструмент миграции. Шаги миграции:

1. **Экспортируйте данные SQLite в SQL**

```bash
sqlite3 data/tokenhub.db .dump > tokenhub_sqlite.sql
```

2. **Конвертируйте синтаксис SQL**

Вручную отредактируйте `tokenhub_sqlite.sql`, адаптировав специфичный для SQLite синтаксис под PostgreSQL.

3. **Импортируйте в PostgreSQL**

```bash
psql -h localhost -U tokenhub -d tokenhub -f tokenhub_sqlite.sql
```

**Примечание**: инструмент для миграции данных запланирован в одной из будущих версий.

## Оптимизация производительности

### Оптимизация индексов

TokenHub автоматически создаёт необходимые индексы, но вы можете добавить дополнительные в зависимости от паттернов запросов:

```sql
-- Если вы часто запрашиваете проекты по cost_center
CREATE INDEX idx_projects_cost_center ON projects(cost_center);

-- Если вы часто запрашиваете записи использования в определённом диапазоне дат
CREATE INDEX idx_usage_records_created_at ON usage_records(created_at);
```

### Мониторинг производительности запросов

Включите журнал медленных запросов PostgreSQL:

```sql
ALTER SYSTEM SET log_min_duration_statement = 1000;  -- Журналировать запросы длиннее 1 секунды
SELECT pg_reload_conf();
```

Просмотр медленных запросов:

```bash
tail -f /var/log/postgresql/postgresql-16-main.log | grep "duration:"
```

## Устранение неполадок

### Сбои подключения

1. **Проверьте, запущен ли PostgreSQL**

```bash
pg_isready -h localhost -U tokenhub
```

2. **Проверьте брандмауэр**

```bash
sudo ufw allow 5432/tcp
```

3. **Проверьте pg_hba.conf**

Убедитесь, что подключения от приложения разрешены:

```
# IPv4 локальные подключения:
host    tokenhub    tokenhub    127.0.0.1/32    md5
```

### Исчерпание пула соединений

При ошибках «too many connections» уменьшите размер пула соединений:

```bash
TOKENHUB_DB_MAX_OPEN_CONNS=10
```

### Проблемы с производительностью

1. **Запустите VACUUM**

```sql
VACUUM ANALYZE;
```

2. **Проверьте раздутие таблиц**

```sql
SELECT schemaname, tablename, pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

## Ссылки

- [Официальная документация PostgreSQL](https://www.postgresql.org/docs/)
- [GORM PostgreSQL Driver](https://github.com/go-gorm/postgres)
- [Документация pgx Driver](https://github.com/jackc/pgx)
