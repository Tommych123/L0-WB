# Демонстрационный сервис с Kafka, PostgreSQL и кешем

## Описание

Демонстрационный сервис на Go для работы с заказами.
Сервис получает данные о заказах из Kafka, сохраняет их в PostgreSQL и кэширует в памяти для быстрого доступа. Также реализован HTTP API и простой веб-интерфейс для просмотра заказов.

---

## Архитектура проекта

```
.
├── cmd/                   # Точка входа приложения
│   └── producer/          # Тестовый скрипт для kafka      
├── internal/
│   ├── config/            # Конфигурации сервиса
│   ├── domain/            # Модели данных
│   ├── http/              # HTTP сервер и обработчики
│   ├── kafka/             # Работа с Kafka
│   ├── repository/        # Работа с PostgreSQL
│   └── service/           # Логика работы с заказами и кеш
├── web/                   # Веб-интерфейс (HTML/JS)
├── migrations/            # SQL миграции
├── docker-compose.yml     # Docker Compose для сервиса, Kafka и PostgreSQL
├── Dockerfile             # Dockerfile приложения
├── go.mod
└── go.sum
```

---

## Требования

* Go 1.24+
* Docker и Docker Compose
* Kafka
* PostgreSQL

---

## Установка и запуск

1. Клонировать репозиторий:

```bash
git clone https://github.com/Tommych123/L0-WB.git
cd L0-WB
```
2. Создать файл с переменными окружения(.env):

```bash
cp .env-example .env
```

3. Запустить сервис, Kafka и PostgreSQL через Docker Compose:

```bash
docker-compose up -d
```

4. Сервис будет доступен по HTTP:

```
http://localhost:8080
```

---

## Использование

### HTTP API

* Получить заказ по `order_uid`:

```
GET http://localhost:8080/order/<order_uid>
```

Возвращает JSON с информацией о заказе. Если заказа нет в кеше, он подтягивается из PostgreSQL.

### Веб-интерфейс

* Ввести `order_uid` в поле ввода и нажать кнопку для получения данных заказа.

---

## Работа с Kafka

* Сервис подписан на топик Kafka и обрабатывает входящие сообщения с заказами.
* Для тестирования можно использовать скрипт-эмулятор отправки сообщений.

---

## База данных

Сервис использует PostgreSQL с четырьмя таблицами:

- **orders** — основная информация о заказе (ID, трек, клиент, дата и др.).  
- **delivery** — данные доставки заказа (имя, адрес, телефон, email), связана с `orders` через `order_uid`.  
- **payment** — информация о платеже (сумма, валюта, банк, дата), связана с `orders` через `order_uid`.  
- **items** — список товаров в заказе (название, цена, количество, бренд), связана с `orders` через `order_uid`.  

---

## Кеширование

* Последние полученные данные заказов хранятся в памяти (map).
* При старте сервиса кеш восстанавливается из базы данных.

---

## Тестирование

Запуск юнит-тестов:

```bash
go test ./internal/service -v
```

Запуск теста для kafka:

```bash
go run cmd/producer/main.go 
```

---

## Модель данных заказа

Пример структуры заказа (`model.json`):

```json
{
   "order_uid": "b563feb7b2b84b6test",
   "track_number": "WBILMTESTTRACK",
   "entry": "WBIL",
   "delivery": {
      "name": "Test Testov",
      "phone": "+9720000000",
      "zip": "2639809",
      "city": "Kiryat Mozkin",
      "address": "Ploshad Mira 15",
      "region": "Kraiot",
      "email": "test@gmail.com"
   },
   "payment": {
      "transaction": "b563feb7b2b84b6test",
      "request_id": "",
      "currency": "USD",
      "provider": "wbpay",
      "amount": 1817,
      "payment_dt": 1637907727,
      "bank": "alpha",
      "delivery_cost": 1500,
      "goods_total": 317,
      "custom_fee": 0
   },
   "items": [
      {
         "chrt_id": 9934930,
         "track_number": "WBILMTESTTRACK",
         "price": 453,
         "rid": "ab4219087a764ae0btest",
         "name": "Mascaras",
         "sale": 30,
         "size": "0",
         "total_price": 317,
         "nm_id": 2389212,
         "brand": "Vivienne Sabo",
         "status": 202
      }
   ],
   "locale": "en",
   "internal_signature": "",
   "customer_id": "test",
   "delivery_service": "meest",
   "shardkey": "9",
   "sm_id": 99,
   "date_created": "2021-11-26T06:22:19Z",
   "oof_shard": "1"
}
```

---

## Примечания

* Некорректные сообщения из Kafka игнорируются или логируются.
* Используются транзакции и подтверждение сообщений для надежности.
* Повторные запросы по одному и тому же `order_uid` обслуживаются из кеша для ускорения.
