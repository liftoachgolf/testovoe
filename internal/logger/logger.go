// logger/logger.go
package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func InitLogger(level string) {
	Logger = logrus.New()

	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		Logger.Fatal("Invalid log level: ", err)
	}
	Logger.SetLevel(logLevel)

	// Настройка формата логирования
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Вывод логов в файл
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Logger.SetOutput(file)
	} else {
		Logger.Info("Failed to log to file, using default stderr")
		file = os.Stderr
	}

	multiWriter := io.MultiWriter(os.Stdout, file)
	Logger.SetOutput(multiWriter)
}
