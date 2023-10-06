package cashe

import (
	"back/pkg/order"
	"errors"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
)

type Cache struct {
	sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	items             map[string]Item
}

type Item struct {
	Value      order.OrderData
	Expiration int64
	Created    time.Time
}

func NewCache(defaultExpiration, cleanupInterval time.Duration, db *sqlx.DB) (*Cache, error) {
	items := make(map[string]Item)
	cache := Cache{
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
		items:             items,
	}

	// Load all orders from the database if the cache is empty
	if len(cache.items) == 0 {
		rows, err := db.Query(`
            SELECT o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature, o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard, d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
                p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee,
                i.chrt_id, i.track_number, i.price, i.rid, i.name, i.sale, i.size, i.total_price, i.nm_id, i.brand, i.status
            FROM orders o
            JOIN delivery d ON o.order_uid = d.order_uid
            JOIN payment p ON o.order_uid = p.order_uid
            JOIN items i ON o.order_uid = i.order_uid
        `)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var ord order.OrderData
			var delivery order.Delivery
			var payment order.Payment
			var item order.Item

			err = rows.Scan(
				&ord.OrderUID,
				&ord.TrackNumber,
				&ord.Entry,
				&ord.Locale,
				&ord.InternalSig,
				&ord.CustomerID,
				&ord.DeliveryService,
				&ord.ShardKey,
				&ord.SmID,
				&ord.DateCreated,
				&ord.OofShard,
				&delivery.Name,
				&delivery.Phone,
				&delivery.Zip,
				&delivery.City,
				&delivery.Address,
				&delivery.Region,
				&delivery.Email,
				&payment.Transaction,
				&payment.RequestID,
				&payment.Currency,
				&payment.Provider,
				&payment.Amount,
				&payment.PaymentDt,
				&payment.Bank,
				&payment.DeliveryCost,
				&payment.GoodsTotal,
				&payment.CustomFee,
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
			)
			if err != nil {
				return nil, err
			}

			ord.Delivery = delivery
			ord.Payment = payment
			ord.Items = append(ord.Items, item)

			// Calculate the remaining time until expiration
			dateCreated, err := time.Parse(time.RFC3339, ord.DateCreated)
			if err != nil {
				return nil, err
			}
			remaining := time.Until(dateCreated.Add(2 * time.Hour))

			// Add the item to the cache if it was created less than two hours ago
			if remaining > 0 {
				cache.Set(ord.OrderUID, ord, remaining)
			}
		}
	}

	if cleanupInterval > 0 {
		cache.StartGC()
	}

	return &cache, nil
}

func (c *Cache) Set(key string, value order.OrderData, duration time.Duration) error {
	var expiration int64

	if duration == 0 {
		duration = c.defaultExpiration
	}

	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	c.Lock()
	defer c.Unlock()

	c.items[key] = Item{
		Value:      value,
		Expiration: expiration,
		Created:    time.Now(),
	}

	return nil
}

func (c *Cache) Get(key string) (*order.OrderData, error) {
	c.RLock()
	defer c.RUnlock()

	item, found := c.items[key]

	// ключ не найден
	if !found {
		return nil, errors.New("key not found")
	}

	// Проверка на установку времени истечения, в противном случае он бессрочный
	if item.Expiration > 0 {
		// Если в момент запроса кеш устарел возвращаем nil
		if time.Now().UnixNano() > item.Expiration {
			return nil, errors.New("cache expired")
		}

	}

	newOrderData := item.Value
	return &newOrderData, nil
}

func (c *Cache) Delete(key string) error {

	c.Lock()

	defer c.Unlock()

	if _, found := c.items[key]; !found {
		return errors.New("key not found")
	}

	delete(c.items, key)

	return nil
}

// Сборка мусора
func (c *Cache) StartGC() {
	go c.GC()
}

func (c *Cache) GC() {
	for {
		<-time.After(c.cleanupInterval)
		if c.items == nil {
			return
		}
		if keys := c.expiredKeys(); len(keys) != 0 {
			c.clearItems(keys)
		}
	}
}

func (c *Cache) expiredKeys() (keys []string) {
	c.RLock()
	defer c.RUnlock()
	for k, i := range c.items {
		if time.Now().UnixNano() > i.Expiration && i.Expiration > 0 {
			keys = append(keys, k)
		}
	}

	return
}

func (c *Cache) clearItems(keys []string) {
	c.Lock()
	defer c.Unlock()
	for _, k := range keys {
		delete(c.items, k)
	}
}
