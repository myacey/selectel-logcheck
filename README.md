# selectel-logcheck

Линтер для проверки корректности лог-сообщений в Go проектах.
Совместим с golangci-lint

## Правила

Линтер проверяет:
1. сообщение начинается со строчной буквы
2. только английский язык
3. отсутствие спецсимволов и эмодзи
4. отсутствие чувствительных данных

## Поддерживаемые логеры

- log
- log/slog
- go.uber.org/zap

## Установка
```sh
go get github.com/myacey/selectel-logcheck
```

## Использование в golangci-lint
Собрать плагин:
```sh
go build -buildmode=plugin -o logcheck.so ./pkg/logcheck/golinters
```

## Конфиг
```yml
linters:
  enable:
    - logcheck

linters-settings:
  logcheck:
    check-lowercase: true
    check-english: true
    check-special: true
    check-sensitive: true

    sensitive-words:
      - username
      - email
```

## Примеры ошибок
- `log.Info("Starting server")` -> `log message should start with lowercase letter`
- `logger.Info("запуск сервера") ` -> `log message should contain only english letters`
- `logger.Info("server started!🚀")` -> `log message should not contain special characters`
- `logger.Info("user password: " + password)` -> `logs should not contain potentially sensitive data`

## Autofix
Линтер поддерживает механизм `SuggestedFix` и может автоматически исправлять часть нарушений.

В данный момент автоматически исправляются:

- приведение первой буквы сообщения к нижнему регистру  
- удаление недопустимых спецсимволов и эмодзи  

### Пример
> [!EXAMPLE]
>  До:
>  ```go
>  logger.Info("Starting server")
>  ```
>  После применения autofix:
>  ```go
>  logger.Info("starting server")
>  ```

> [!EXAMPLE]
>  До:
>  ```go
>  logger.Info("!!!starting server!!!")
>  ```
>  После применения autofix:
>  ```go
>  logger.Info("starting server")
>  ```
