# Wallet API (Golang + PostgreSQL)

## 📌 Описание
Приложение для управления балансом кошельков через REST API.  
Поддерживает операции `DEPOSIT` и `WITHDRAW`, а также получение баланса.

## 🚀 Запуск
### 1. Клонировать репозиторий
```bash
git clone https://github.com/ТВОЙ_РЕПОЗИТОРИЙ.git
cd wallet_go
Запустить Контейнеры 
docker-compose up --build -d
Тест API с curl
curl -X POST "http://localhost:8080/api/v1/wallet" \
     -H "Content-Type: application/json" \
     -d '{"walletId":"550e8400-e29b-41d4-a716-446655440000","operationType":"DEPOSIT","amount":1000}'
Проверить баланс:
curl -X GET "http://localhost:8080/api/v1/wallets/550e8400-e29b-41d4-a716-446655440000"
Запустить нагрузочный тест (Artillery)
artillery run artillery_test.yml
Стек технологий
Golang
PostgreSQL
Docker + Docker Compose
Gin (Web API)
Logrus (логирование)
Artillery (нагрузочное тестирование)
Unit-тесты (Testify)
