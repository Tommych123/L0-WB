package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Tommych123/L0-WB/internal/domain"
	"github.com/jmoiron/sqlx"
)

// Описание структуры для подключения к БД
type PostgresOrderRepository struct {
	db *sqlx.DB
}

// Создание нового репозитория
func NewPostgresOrderRepository(db *sqlx.DB) *PostgresOrderRepository {
	return &PostgresOrderRepository{db: db}
}

// Save сохраняет заказ в БД
func (rep *PostgresOrderRepository) Save(order *domain.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := rep.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	// Вставка в orders
	_, err = tx.ExecContext(ctx, `
	INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature,
		customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		ON CONFLICT (order_uid) DO NOTHING`,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale,
		order.InternalSignature, order.CustomerID, order.DeliveryService,
		order.ShardKey, order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		tx.Rollback()
		return err
	}
	// Вставка в delivery
	_, err = tx.ExecContext(ctx, `
        INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
        ON CONFLICT (order_uid) DO NOTHING
    `,
		order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Вставка в payment
	_, err = tx.ExecContext(ctx, `
        INSERT INTO payment (order_uid, transaction, request_id, currency, provider, amount,
                             payment_dt, bank, delivery_cost, goods_total, custom_fee)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
        ON CONFLICT (order_uid) DO NOTHING
    `,
		order.OrderUID, order.Payment.Transaction, order.Payment.RequestID,
		order.Payment.Currency, order.Payment.Provider, order.Payment.Amount,
		order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost,
		order.Payment.GoodsTotal, order.Payment.CustomFee,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Вставка в items (может быть несколько записей)
	for _, item := range order.Items {
		_, err = tx.ExecContext(ctx, `
            INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size,
                               total_price, nm_id, brand, status, order_uid)
            VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
            ON CONFLICT (rid) DO NOTHING
        `,
			item.ChrtID, item.TrackNumber, item.Price, item.RID, item.Name,
			item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand,
			item.Status, order.OrderUID,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// Get получение записи по ID
func (r *PostgresOrderRepository) Get(orderUID string) (*domain.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var order domain.Order
	err := r.db.GetContext(ctx, &order, `
        SELECT order_uid, track_number, entry, locale, internal_signature,
               customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
        FROM orders
        WHERE order_uid = $1
    `, orderUID)
	if err != nil {
		// Заказ отсутствует
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	// Получаем delivery
	var delivery domain.Delivery
	err = r.db.GetContext(ctx, &delivery, `
        SELECT name, phone, zip, city, address, region, email
        FROM delivery WHERE order_uid = $1
    `, orderUID)
	if err != nil {
		return nil, err
	}
	order.Delivery = delivery

	// Получаем payment
	var payment domain.Payment
	err = r.db.GetContext(ctx, &payment, `
        SELECT transaction, request_id, currency, provider, amount,
               payment_dt, bank, delivery_cost, goods_total, custom_fee
        FROM payment WHERE order_uid = $1
    `, orderUID)
	if err != nil {
		return nil, err
	}
	order.Payment = payment

	// Получаем items
	var items []domain.Item
	err = r.db.SelectContext(ctx, &items, `
        SELECT chrt_id, track_number, price, rid, name, sale, size,
               total_price, nm_id, brand, status
        FROM items WHERE order_uid = $1
    `, orderUID)
	if err != nil {
		return nil, err
	}
	order.Items = items

	// Возвращаем собранный заказ из 4 таблиц
	return &order, nil
}
