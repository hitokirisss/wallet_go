package main

import (
	"log"

	"wallet_go/internal/database"
	"wallet_go/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные окружения
	if err := godotenv.Load("config.env"); err != nil {
		log.Fatal("Ошибка загрузки config.env")
	}

	// Инициализируем БД
	database.InitDB()

	// Запускаем сервер
	srv := server.NewServer(database.DB)
	srv.Run("0.0.0.0:8080")
}
