package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/1337yeeee/order-service-wb/internal/cache"
	"github.com/1337yeeee/order-service-wb/internal/models"
)

type OrderRepository struct {
	cache *cache.Cache
	db    *sql.DB
}

func New(cache *cache.Cache, db *sql.DB) *OrderRepository {
	return &OrderRepository{
		cache: cache,
		db:    db,
	}
}

func (r *OrderRepository) Save(order models.Order) error {
	// Сохраняем заказ в кэш
	if err := r.saveToCache(order); err != nil {
		return err
	}

	// Сохраняем заказ в базу данных
	if err := r.saveToDB(order); err != nil {
		return fmt.Errorf("failed to save order to DB: %w", err)
	}

	return nil
}

func (r *OrderRepository) saveToCache(order models.Order) error {
	// Сериализируем модель заказа
	orderData, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %w", err)
	}

	// Сохраняем заказ в кэш
	r.cache.Set(order.OrderUID, orderData)
	log.Printf("Order cached (ID: %s)", order.OrderUID)

	return nil
}

func (r *OrderRepository) saveToDB(order models.Order) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Сохраняем основной заказ
	_, err = tx.Exec(`
		INSERT INTO orders (
			order_uid, track_number, entry, locale, internal_signature,
			customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (order_uid) DO NOTHING`,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.Shardkey,
		order.SmID,
		order.DateCreated,
		order.OofShard,
	)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	// Сохраняем доставку
	_, err = tx.Exec(`
		INSERT INTO deliveries (
			order_uid, name, phone, zip, city, address, region, email
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (order_uid) DO NOTHING`,
		order.OrderUID,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
	)
	if err != nil {
		return fmt.Errorf("failed to insert delivery: %w", err)
	}

	// Сохраняем платеж
	_, err = tx.Exec(`
		INSERT INTO payments (
			order_uid, transaction, request_id, currency, provider,
			amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (order_uid) DO NOTHING`,
		order.OrderUID,
		order.Payment.Transaction,
		order.Payment.RequestID,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDt,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee,
	)
	if err != nil {
		return fmt.Errorf("failed to insert payment: %w", err)
	}

	// Сохраняем товары
	for _, item := range order.Items {
		_, err = tx.Exec(`
			INSERT INTO items (
				order_uid, chrt_id, track_number, price, rid, name,
				sale, size, total_price, nm_id, brand, status
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
			order.OrderUID,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.Rid,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status,
		)
		if err != nil {
			return fmt.Errorf("failed to insert item: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Order saved in DB (ID: %s)", order.OrderUID)

	return nil
}

func (r *OrderRepository) Get(orderUID string) (*models.Order, error) {
	// Пробуем получить из кэша
	order, err := r.getFromCache(orderUID)
	if order != nil || err != nil {
		return order, err
	}

	// Если нет в кэше, ищем в БД
	order, err = r.getFromDB(orderUID)
	if err == nil {
		log.Printf("Order found in DB (ID: %s)", order.OrderUID)
	}

	r.saveToCache(*order)

	return order, err
}

func (r *OrderRepository) getFromCache(orderUID string) (*models.Order, error) {
	data, ok := r.cache.Get(orderUID)
	if !ok {
		return nil, nil
	}

	var order models.Order
	if err := json.Unmarshal(data, &order); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached order: %w", err)
	}
	log.Printf("Order found in cache (ID: %s)", order.OrderUID)

	return &order, nil
}

func (r *OrderRepository) getFromDB(orderUID string) (*models.Order, error) {
	// Получаем основную информацию о заказе
	var order models.Order
	err := r.db.QueryRow(`
		SELECT order_uid, track_number, entry, locale, 
		       internal_signature, customer_id, delivery_service,
		       shardkey, sm_id, date_created, oof_shard
		FROM orders WHERE order_uid = $1`, orderUID).Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.Shardkey,
		&order.SmID,
		&order.DateCreated,
		&order.OofShard,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Получаем delivery
	err = r.db.QueryRow(`
		SELECT name, phone, zip, city, address, region, email
		FROM deliveries WHERE order_uid = $1`, orderUID).Scan(
		&order.Delivery.Name,
		&order.Delivery.Phone,
		&order.Delivery.Zip,
		&order.Delivery.City,
		&order.Delivery.Address,
		&order.Delivery.Region,
		&order.Delivery.Email,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get delivery: %w", err)
	}

	// Получаем payment
	err = r.db.QueryRow(`
		SELECT transaction, request_id, currency, provider,
		       amount, payment_dt, bank, delivery_cost,
		       goods_total, custom_fee
		FROM payments WHERE order_uid = $1`, orderUID).Scan(
		&order.Payment.Transaction,
		&order.Payment.RequestID,
		&order.Payment.Currency,
		&order.Payment.Provider,
		&order.Payment.Amount,
		&order.Payment.PaymentDt,
		&order.Payment.Bank,
		&order.Payment.DeliveryCost,
		&order.Payment.GoodsTotal,
		&order.Payment.CustomFee,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	// Получаем items
	rows, err := r.db.Query(`
		SELECT chrt_id, track_number, price, rid, name,
		       sale, size, total_price, nm_id, brand, status
		FROM items WHERE order_uid = $1`, orderUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item models.Item
		if err := rows.Scan(
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
		); err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		order.Items = append(order.Items, item)
	}

	return &order, nil
}
