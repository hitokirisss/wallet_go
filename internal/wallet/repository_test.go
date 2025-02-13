package wallet

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var testDB *sql.DB
var repo *Repository

func TestMain(m *testing.M) {
	// Подключаемся к тестовой базе
	var err error
	testDB, err = sql.Open("postgres", "postgres://postgres:123123@localhost:5432/test_db?sslmode=disable")
	if err != nil {
		log.Fatalf("❌ Ошибка подключения к тестовой базе: %v", err)
	}

	repo = NewRepository(testDB)

	// Выполняем тесты
	code := m.Run()

	// Закрываем соединение после тестов
	testDB.Close()

	// Выходим с кодом выполнения тестов
	os.Exit(code)
}

func TestGetWalletBalance(t *testing.T) {
	// Создаем тестовый кошелек
	walletID := uuid.New()
	_, err := testDB.Exec("INSERT INTO wallets (id, balance) VALUES ($1, $2)", walletID, 1000)
	assert.NoError(t, err)

	// Запрашиваем баланс
	balance, err := repo.GetWalletBalance(walletID)
	assert.NoError(t, err)
	assert.Equal(t, int64(1000), balance, "Баланс кошелька должен быть 1000")

	// Удаляем тестовые данные
	_, err = testDB.Exec("DELETE FROM wallets WHERE id = $1", walletID)
	assert.NoError(t, err)
}
