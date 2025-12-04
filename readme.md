# Wallet Service

Мини-проектория: сервис кошельков с базовыми операциями.  
Пополнение, списание, получение текущего баланса.  
Хранение в PostgreSQL, запуск через Docker Compose.

## Как запустить

1. Разместить файл `.env` рядом с `docker-compose.yml`  

2. Поднять сервис:
docker compose up --build

После сборки приложение будет доступно на `http://localhost:8080`.

## REST API

### POST /api/v1/wallet  
Операция над кошельком.

Тело:
{
"walletId": "UUID",
"operationType": "DEPOSIT" | "WITHDRAW",
"amount": 1000
}


Пример ответа:
{
"walletId": "UUID",
"balance": 11000
}


---

### GET /api/v1/wallets/{walletId}  
Получение баланса.

Ответ:
{
"walletId": "UUID",
"balance": 5000
}


## Тесты

Запуск локально:
go test ./... -cover

Проект покрыт тестами на конкурентную работу.  

конкурентная безопасность обеспечивается тем, что update wallets set balance = balance + delta выполняется атомарно под каоптом постгри и блокирует строку на время операции. из-за  row-level locking и проверке что balance + delta >= 0 в том же апдейте, несколько параллельных запросов не могут повредить баланс.