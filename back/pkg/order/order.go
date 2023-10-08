package order

import (
	"github.com/jmoiron/sqlx"
)

type OrderData struct {
	OrderUID        string   `json:"order_uid"`
	TrackNumber     string   `json:"track_number"`
	Entry           string   `json:"entry"`
	Locale          string   `json:"locale"`
	InternalSig     string   `json:"internal_signature"`
	CustomerID      string   `json:"customer_id"`
	DeliveryService string   `json:"delivery_service"`
	ShardKey        string   `json:"shardkey"`
	SmID            int      `json:"sm_id"`
	DateCreated     string   `json:"date_created"`
	OofShard        string   `json:"oof_shard"`
	Delivery        Delivery `json:"delivery"`
	Payment         Payment  `json:"payment"`
	Items           []Item   `json:"items"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int64  `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Item struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

// Вставка данных в таблицу "orders"
func InsertOrder(db *sqlx.DB, data OrderData) error {
	_, err := db.Exec(`
        INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		data.OrderUID, data.TrackNumber, data.Entry, data.Locale, data.InternalSig, data.CustomerID, data.DeliveryService, data.ShardKey, data.SmID, data.DateCreated, data.OofShard)
	return err
}

// Вставка данных в таблицу "delivery"
func InsertDelivery(db *sqlx.DB, data OrderData) error {
	_, err := db.Exec(`
        INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		data.OrderUID, data.Delivery.Name, data.Delivery.Phone, data.Delivery.Zip, data.Delivery.City, data.Delivery.Address, data.Delivery.Region, data.Delivery.Email)
	return err
}

// Вставка данных в таблицу "payment"
func InsertPayment(db *sqlx.DB, data OrderData) error {
	_, err := db.Exec(`
        INSERT INTO payment (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		data.OrderUID, data.Payment.Transaction, data.Payment.RequestID, data.Payment.Currency, data.Payment.Provider, data.Payment.Amount, data.Payment.PaymentDt, data.Payment.Bank, data.Payment.DeliveryCost, data.Payment.GoodsTotal, data.Payment.CustomFee)
	return err
}

// Вставка данных в таблицу "items"
func InsertItems(db *sqlx.DB, data OrderData) error {
	for _, item := range data.Items {
		_, err := db.Exec(`
            INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
			data.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			return err
		}
	}
	return nil
}
