package wallet

import (
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"wallet_go/internal/logger"
)

type Repository struct {
	db *sql.DB
}

// Создание нового репозитория
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Получить баланс кошелька
func (r *Repository) GetWalletBalance(walletID uuid.UUID) (int64, error) {
	var balance int64
	// Запрос к базе данных для получения баланса кошелька
	err := r.db.QueryRow("SELECT balance FROM wallets WHERE id = $1", walletID).Scan(&balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("кошелек не найден")
		}
		log.Printf("Ошибка при чтении баланса: %v", err)
		return 0, err
	}
	return balance, nil
}

// Функция обновления баланса кошелька
func (r *Repository) UpdateBalance(walletID uuid.UUID, amount int64, operationType string) error {
	log := logger.GetLogger()
	tx, err := r.db.Begin()
	if err != nil {
		log.WithError(err).Error("Ошибка начала транзакции")
		return err
	}
	defer tx.Rollback()

	// Устанавливаем уровень изоляции транзакции
	_, err = tx.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
	if err != nil {
		log.WithError(err).Error("Ошибка установки уровня изоляции")
		return err
	}

	log.WithFields(logrus.Fields{"walletId": walletID}).Info("🔄 Начало обновления баланса")

	// 🔥 Блокируем строку с балансом для обновления
	var currentBalance int64
	err = tx.QueryRow("SELECT balance FROM wallets WHERE id = $1 FOR UPDATE", walletID).Scan(&currentBalance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn("Кошелек не найден")
			return errors.New("кошелек не найден")
		}
		log.WithError(err).Error("Ошибка при чтении данных из БД")
		return err
	}

	// Рассчитываем новый баланс
	var newBalance int64
	if operationType == "DEPOSIT" {
		newBalance = currentBalance + amount
	} else if operationType == "WITHDRAW" {
		if currentBalance < amount {
			log.Warn("Недостаточно средств для вывода")
			return errors.New("недостаточно средств")
		}
		newBalance = currentBalance - amount
	} else {
		log.Warn("Неизвестный тип операции")
		return errors.New("неизвестный тип операции")
	}

	// Обновляем баланс в таблице wallets
	_, err = tx.Exec("UPDATE wallets SET balance = $1 WHERE id = $2", newBalance, walletID)
	if err != nil {
		log.WithError(err).Error("Ошибка при обновлении баланса")
		return err
	}

	// 📌 Вставляем запись в transactions
	_, err = tx.Exec("INSERT INTO transactions (wallet_id, operation_type, amount) VALUES ($1, $2, $3)",
		walletID, operationType, amount)
	if err != nil {
		log.WithError(err).Error("Ошибка при записи транзакции")
		return err
	}

	// Фиксируем транзакцию
	err = tx.Commit()
	if err != nil {
		log.WithError(err).Error("Ошибка при коммите транзакции")
		return err
	}

	log.WithFields(logrus.Fields{"walletId": walletID, "newBalance": newBalance}).Info("✅ Баланс обновлен успешно")
	return nil
}
