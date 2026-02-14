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

## Установка как Go-анализатора
```go
go install github.com/myacey/selectel-logcheck/cmd/addlint@latest
```
### Запуск:
```sh
addlint ./...
```

## Использование в golangci-lint
### 1. Скачать исходники golangci-lint
```sh
git clone https://github.com/golangci/golangci-lint
cd golangci-lint
```

### 2. Собрать бинарник
```sh
make build
```
Готовый файл:
```sh
./build/golangci-lint
```

### 3. Собрать плагин
В проекте selectel-logcheck:
```
go build -buildmode=plugin -o logcheck.so ./pkg/logcheck/golinters
```

### 4. Настроить `golangci.yml`
```yml
version: "2"

linters:
  enable:
    - logcheck

  settings:
    custom:
      logcheck:
        path: ./logcheck.so
        settings:
          check-lowercase: true
          check-english: true
          check-special: true
          check-sensitive: true
          sensitive-words:
            - username
            - email
```

### 5. Запускать кастомный golangci-lint
```sh
/path/to/golangci-lint/build/golangci-lint run ./...
```
> [!WARN]
> Не глобальный из `$PATH`.

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
> [!NOTE]
>  До:
>  ```go
>  logger.Info("Starting server")
>  ```
>  После применения autofix:
>  ```go
>  logger.Info("starting server")
>  ```

> [!NOTE]
>  До:
>  ```go
>  logger.Info("!!!starting server!!!")
>  ```
>  После применения autofix:
>  ```go
>  logger.Info("starting server")
>  ```
