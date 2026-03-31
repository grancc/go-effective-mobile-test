# Go Effectiva Mobile Test

REST API на Go для данных о подписках

## Запуск через Docker Compose

Из корня проекта:

```bash
docker compose up --build
```

Сервисы:

- **db** — PostgreSQL, порт `5432` на хосте
- **migrate** — применяет SQL-миграции из каталога `shema/`
- **web** — API на `http://localhost:8080`

После старта откройте Swagger: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html).

## Локальный запуск

1. Поднимите PostgreSQL и задайте `POSTGRES_*` в окружении.

2. В `configs/config.yaml` укажите `db.host: localhost` (и при необходимости `db.port`), если база не в Docker-сети.

3. Примените миграции (пример с установленным `migrate`):

   ```bash
   migrate -path shema -database "postgres://USER:PASSWORD@localhost:5432/DBNAME?sslmode=disable" up
   ```

4. Запуск из корня репозитория:

   ```bash
   go run ./cmd/main.go
   ```


## Структура проекта (кратко)

- `cmd/main.go` — точка входа
- `pkg/handler` — HTTP (Gin)
- `pkg/service` — бизнес-логика
- `pkg/repository` — PostgreSQL (sqlx)
- `shema/` — SQL-миграции (имя каталога сохранено как в репозитории)
- `docs/` — сгенерированные файлы Swagger
