# Eshop Seller Service

Микросервис продавцов для проекта **Eshop**.

Сервис хранит данные продавцов и их внешние ссылки, предоставляет gRPC API и использует PostgreSQL в качестве основной базы данных.

## Возможности

### Seller

- создание продавца;
- получение продавца по идентификатору;
- получение продавцов пользователя;
- получение статуса продавца;
- обновление данных продавца;
- архивирование продавца;
- удаление продавца.

### Social Link

- добавление ссылки продавца;
- получение списка ссылок продавца;
- удаление ссылки.

Обновление ссылки отдельным RPC не предусмотрено. Некорректная ссылка удаляется и создаётся заново.

## Технологии

- Go 1.26;
- gRPC;
- Protocol Buffers;
- PostgreSQL;
- pgx;
- UUID;
- cleanenv;
- Taskfile;
- SQL-миграции.

Контракты gRPC находятся в отдельном модуле:

```text
github.com/barnigator/protos
```

## Архитектура

Сервис разделён на несколько слоёв:

```text
gRPC handler
      ↓
   usecase
      ↓
   domain
      ↓
PostgreSQL repository
```

Основные каталоги:

```text
cmd/seller-service          точка входа приложения
config                      конфигурационные файлы
examples/seller-client      клиент для ручной проверки gRPC API
internal/domain             доменные сущности и ошибки
internal/usecase            бизнес-логика
internal/repository         работа с PostgreSQL
internal/grpc/handler       gRPC handlers
internal/grpc/server        запуск gRPC-сервера
migrations                  SQL-миграции
```

## Требования

Для локального запуска необходимы:

- Go;
- PostgreSQL;
- Task;
- применённые SQL-миграции из каталога `migrations`.

## Конфигурация

Локальная конфигурация находится в файле:

```text
config/local.yaml
```

Пример:

```yaml
env: "local"

app:
  startup_timeout: 5s

grpc:
  port: 44045
  timeout: 5s

postgres:
  dsn: ""
```

Строка подключения к PostgreSQL передаётся через переменную окружения:

```env
POSTGRES_DSN=postgres://eshop_seller:eshop_seller@localhost:5432/eshop_seller?sslmode=disable
```

Создай файл `.env` на основе `.env.example`:

```bash
cp .env.example .env
```

На Windows файл можно скопировать вручную и заполнить актуальными параметрами подключения.

Путь к конфигурации передаётся флагом:

```bash
--config=./config/local.yaml
```

или через переменную окружения:

```env
CONFIG_PATH=./config/local.yaml
```

## Миграции

SQL-миграции находятся в каталоге:

```text
migrations
```

Перед запуском сервиса миграции должны быть применены к базе данных.

В общем окружении проекта миграции рекомендуется запускать через репозиторий `eshop-infra`.

## Установка зависимостей

```bash
go mod download
```

Для очистки и проверки зависимостей:

```bash
go mod tidy
```

## Локальный запуск

Запуск через Taskfile:

```bash
task run
```

Эквивалентная команда:

```bash
go run ./cmd/seller-service --config=./config/local.yaml
```

По умолчанию gRPC-сервер запускается на порту:

```text
44045
```

## Ручная проверка

После запуска сервиса можно запустить тестовый gRPC-клиент:

```bash
task client
```

или:

```bash
go run ./examples/seller-client
```

Клиент выполняет положительные и отрицательные сценарии для Seller и Social Link.

Некоторые вызовы намеренно возвращают ошибки, например:

- `InvalidArgument` для пустого или некорректного UUID;
- `NotFound` для отсутствующего продавца или ссылки;
- повторное удаление ссылки возвращает `NotFound`.

Такие ответы являются ожидаемой частью smoke-проверки.

## Сборка

```bash
go build ./...
```

Сборка исполняемого файла сервиса:

```bash
go build -o seller-service ./cmd/seller-service
```

## Тесты

Тесты временно не являются обязательным этапом MVP, но Taskfile уже содержит команды:

```bash
task test
task test-usecase
task test-integration
```

Интеграционные тесты используют переменную:

```env
POSTGRES_TEST_DSN=postgres://eshop_seller:eshop_seller@localhost:5432/eshop_seller_test?sslmode=disable
```

## Основные команды Taskfile

```bash
task                 показать доступные команды
task run             запустить сервис
task client          запустить клиент
task test            запустить все тесты
task test-usecase    запустить тесты usecase
task test-integration запустить интеграционные тесты PostgreSQL
```

## Текущее состояние

Для MVP реализованы:

- Seller CRUD;
- статус продавца;
- архивирование продавца;
- Social Link: Add, List, Delete;
- PostgreSQL repository;
- gRPC API;
- SQL-миграции;
- клиент для ручной проверки.

Следующий этап — контейнеризация сервиса и подключение его к общему окружению `eshop-infra`.
