package p

import (
	"fmt"
	"log"
	"log/slog"

	"go.uber.org/zap"
)

func example() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	//❌Неправильно
	logger.Info("Starting server on port 8080") // want "log message should start with lowercase letter"
	slog.Error("Failed to connect to database") // want "log message should start with lowercase letter"
	//✅Правильно
	logger.Info("starting server on port 8080") // ok
	slog.Error("failed to connect to database") // ok

	// -----------------------------

	//❌Неправильно
	logger.Info("запуск сервера")                    // want "log message should contain only english letters"
	logger.Error("ошибка подключения к базе данных") // want "log message should contain only english letters"
	//✅Правильно
	logger.Info("starting server")                // ok
	logger.Error("failed to connect to database") // ok

	// -----------------------------

	//❌Неправильно
	logger.Info("server started!🚀")                 // want "log message should not contain special characters"
	logger.Error("connection failed!!!")            // want  "log message should not contain special characters"
	logger.Warn("warning: something went wrong...") // want "log message should not contain special characters"
	//✅Правильно
	logger.Info("server started")       // ok
	logger.Error("connection failed")   // ok
	logger.Warn("something went wrong") // ok

	// -----------------------------

	//❌Неправильно
	password := "a"
	apiKey := "b"
	token := "c"
	logger.Info("user password: " + password) // want "logs should not contain potentially sensitive data"
	logger.Debug("api_key=" + apiKey)         // want "logs should not contain potentially sensitive data"
	logger.Info("token: " + token)            // want "logs should not contain potentially sensitive data"
	//✅Правильно
	logger.Info("user authenticated successfully") // ok
	logger.Debug("api request completed")          // ok
	logger.Info("token validated")                 // ok
}

func checkLog() {
	log.Print("hello World")
	log.Println("Hello World") // want "log message should start with lowercase letter"
	log.Panic("hello мир")     // want "log message should contain only english letters"
	log.Print("hello world!")  // want "log message should not contain special characters"

	password := "a"
	log.Print("password: " + password) // want "logs should not contain potentially sensitive data"
}

func checkSlog() {
	slog.Info("hello World")
	slog.Warn("Hello World")  // want "log message should start with lowercase letter"
	slog.Error("hello мир")   // want "log message should contain only english letters"
	slog.Info("hello world!") // want "log message should not contain special characters"

	password := "a"
	slog.Debug("password: " + password) // want "logs should not contain potentially sensitive data"
}

func checkZap() {
	zap.L().Info("hello World")
	zap.L().Warn("Hello World")   // want "log message should start with lowercase letter"
	zap.L().Error("hello мир")    // want "log message should contain only english letters"
	zap.L().Debug("hello world!") // want "log message should not contain special characters"

	password := "a"
	zap.L().Info("password: " + password) // want "logs should not contain potentially sensitive data"
}

func sensitiveEdgecases() {
	pass := "x"
	zap.L().Info("user password: " + pass + "!!!") // want "logs should not contain potentially sensitive data"

	userToken := "abc"
	zap.L().Info("validated", userToken) // want "logs should not contain potentially sensitive data"

	password := "x"
	zap.L().Info("password is", password) // want "logs should not contain potentially sensitive data"

	zap.L().Info(fmt.Sprintf(password)) // want "logs should not contain potentially sensitive data"
}

func safeSensitiveWords() {
	name := "bob"
	zap.L().Info("hello " + name) // ok

	zap.L().Info("password validated") // ok
}
