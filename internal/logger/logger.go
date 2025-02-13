package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func InitLogger() {
	log.SetFormatter(&logrus.JSONFormatter{}) // JSON-логирование
	log.SetOutput(os.Stdout)                  // Вывод в консоль
	log.SetLevel(logrus.DebugLevel)           // Уровень логирования
}

// Функция для получения логгера
func GetLogger() *logrus.Logger {
	return log
}
