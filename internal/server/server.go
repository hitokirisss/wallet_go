package server

import (
	"database/sql"
	"log"
	"net/http"
	"wallet_go/internal/logger"

	"wallet_go/internal/wallet"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Server struct {
	router     *gin.Engine
	repository *wallet.Repository
}

func NewServer(db *sql.DB) *Server {
	logger.InitLogger() // Инициализация логирования

	r := gin.Default()
	repo := wallet.NewRepository(db)

	s := &Server{
		router:     r,
		repository: repo,
	}

	s.routes()
	logger.GetLogger().Info("🚀 Сервер успешно инициализирован!")

	return s
}

func (s *Server) routes() {
	s.router.POST("/api/v1/wallet", s.handleWalletTransaction)          // Создание и операция с кошельком
	s.router.GET("/api/v1/wallets/:walletID", s.handleGetWalletBalance) // Получение баланса кошелька
}

func (s *Server) handleWalletTransaction(c *gin.Context) {
	log := logger.GetLogger()

	var req struct {
		WalletID      uuid.UUID `json:"walletId" binding:"required"`
		OperationType string    `json:"operationType" binding:"required"`
		Amount        int64     `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithError(err).Warn("❌ Неверный формат запроса")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	log.WithFields(logrus.Fields{
		"walletId":      req.WalletID,
		"operationType": req.OperationType,
		"amount":        req.Amount,
	}).Info("📩 Получен новый запрос")

	err := s.repository.UpdateBalance(req.WalletID, req.Amount, req.OperationType)
	if err != nil {
		log.WithError(err).Error("❌ Ошибка обработки запроса")
		if err.Error() == "недостаточно средств" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Внутренняя ошибка сервера"})
		}
		return
	}

	log.WithFields(logrus.Fields{"walletId": req.WalletID}).Info("✅ Операция выполнена успешно")
	c.JSON(http.StatusOK, gin.H{"message": "Операция выполнена успешно"})
}

func (s *Server) handleGetWalletBalance(c *gin.Context) {
	walletID, err := uuid.Parse(c.Param("walletID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат UUID"})
		return
	}

	balance, err := s.repository.GetWalletBalance(walletID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Кошелек не найден"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"walletId": walletID, "balance": balance})
}

func (s *Server) Run(addr string) {
	log.Printf("🚀 Сервер запущен на %s", addr)
	if err := s.router.Run(addr); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
