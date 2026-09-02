# Архитектура фреймворка миграции

## Обзор

Фреймворк миграции использует трёхфазную архитектуру:

1. **Source Adapter** — считывает конфигурацию конкурирующего шлюза и формирует `CanonicalMigrationBundle`
2. **Canonical Bundle** — версионированное промежуточное представление JSON без секретов в открытом виде
3. **TokenHub Sink** — идемпотентно применяет bundle к TokenHub через Admin API

## Добавление нового адаптера источника

1. Реализуйте интерфейс `source.Extractor` в `backend/internal/migration/source/<name>/`
2. Зарегистрируйте адаптер через `init()` с помощью `source.Register()`
3. Добавьте фикстуры в `testdata/`
4. Добавьте документацию в `docs/migration/<name>.md`

### Интерфейс Extractor

```go
type Extractor interface {
    Name() string
    SupportedVersions() []string
    Probe(ctx context.Context, opts ExtractOptions) (Info, error)
    Extract(ctx context.Context, opts ExtractOptions) (*bundle.CanonicalMigrationBundle, error)
}
```

## Операции Sink

- **Plan** — пробный запуск (dry-run), сообщает о планируемых изменениях
- **Apply** — идемпотентный upsert в TokenHub
- **Verify** — проверяет соответствие применённого состояния bundle
- **Rollback** — откатывает состояние до применения с использованием контрольной точки (checkpoint)
