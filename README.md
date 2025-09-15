# Сервис регистрации ошибок и отправкой их в Sentry

## Описание

- **Catcher** — сервис на Go для регистрации ошибок и отправки их в Sentry.
- Главная точка входа: `app/cmd/catcher/main.go`.
- Основная логика находится в `internal/` (config, git, handler, lib, models, sentryhub, server, service, testutil).
- Конфигурационные файлы — в `config/`.
- Скрипты для Windows-сервиса — `service/windows/`.
- Swagger/OpenAPI документация — `docs/` и файлы `swagger.*`.
  
## Основные методы HTTP API

- `POST /api/reg/pushReport` — регистрация и отправка отчёта об ошибке в Sentry. Принимает структуру отчёта, возвращает идентификатор события.
- `POST /api/prj/:id/sendEvent` — ручная отправка события в Sentry. Позволяет отправить кастомное событие с произвольными параметрами. В теле запроса указываются необходимые поля события.
- `POST /api/service/clearCache` — очистка кэша (например, для сброса кешированных исключений или контекста).
- `GET /swagger/*any` — Swagger UI и документация для всех методов сервиса.

## Сборка

- Сборка для Windows: `make build-win` или `go build .\app\cmd\catcher`.
- Swagger-документация: `swag init -g app/cmd/catcher/main.go --parseDependency`.

## Варианты запуска

Сервис поддерживает запуск с различными параметрами командной строки:

- `-config <путь>` — указать путь к конфигурационному файлу (по умолчанию `config/config.yml`)
- `-debug` — запустить сервис в режиме отладки (использует `config_debug.yml`)
- `-port <номер>` — задать порт сервера вручную

Примеры:

```powershell
- go run .\app\cmd\catcher -config config/config_debug.yml -port 8080 -debug
- catcher.exe -config config/config.yml
```

## Структура файла config.yml

### Server

- `Port`: порт, на котором запускается сервис (например, "8000").

### Registry

- `UserMessage`: текст для пользователя при ошибке.
- `DumpType`: тип дампа для отчёта.
- `Timeout`: таймаут ожидания ответа.

### Projects

- Список проектов, для которых работает сервис.
- В каждом проекте:
  - `Name`, `ID`: имя и идентификатор проекта.
  - **Service**: отвечает за параметры интеграции с сервисом для получения дополнительных данных:
    - `Use`: включить/отключить интеграцию с сервисом (true/false)
    - `Url`: адрес внешнего сервиса, куда отправляются запросы
    - `IimeOut`: таймаут ожидания ответа от сервиса (в минутах)
    - `Credintials`:
      - `UserName`: имя пользователя для авторизации
      - `Password`: пароль пользователя
    - `Cache`:
      - `Use`: использовать кеширование запросов к сервису
      - `Expiration`: время жизни кеша (в минутах)
    - `Exeptions`:
      - `Use`: использовать обработку исключений
      - `Cache`:
        - `Use`: кешировать исключения
        - `Expiration`: время жизни кеша исключений (в минутах)
    - `Test`:
      - `UserName`: имя тестового пользователя для проверки интеграции
  - **Sentry**: интеграция с Sentry
    - `Dsn`: DSN для Sentry
    - `Environment`: окружение (prod/dev)
    - `Platform`: платформа
    - `ContextAround`: параметры контекста ошибки
    - `Attachments`: вложения к ошибке
    - `SendingCache`: кеш отправки
  - **Git**: параметры доступа к исходному коду
    - `Use`: включить интеграцию
    - `Url`, `Path`, `Token`, `Branch`, `SourceCodeRoot`: настройки репозитория
  - **Extentions**: список расширений

### Log

- `Debug`: режим отладки
- `Level`: уровень логирования (2-5)
- `OutputInFile`: логирование в файл
- `Dir`: каталог логов

### DeleteTempFiles

- `true`: удалять временные файлы из `web/temp`

### Sentry (внутренний)

- `Use`: включить Sentry для внутренних ошибок
- `Dsn`: DSN для Sentry
- `AttachStacktrace`: прикладывать stacktrace
- `TracesSampleRate`: частота трассировки
- `EnableTracing`: включить трассировку

## TODO
