package repository

import (
	"back/pkg/order"

	"github.com/jmoiron/sqlx"
)

type GetOrderByIdPostgres struct {
	db *sqlx.DB
}

func NewGetOrderByIdPostgres(db *sqlx.DB) *GetOrderByIdPostgres {
	return &GetOrderByIdPostgres{db: db}
}

func (r *GetOrderByIdPostgres) GetById(id string) (*order.OrderData, error) {
	var o order.OrderData
	err := r.db.Get(&o, `
        SELECT o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature, o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard,
               d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
               p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee,
               i.chrt_id, i.price, i.rid, i.name, i.sale, i.size, i.total_price, i.nm_id, i.brand, i.status
        FROM orders o
        JOIN delivery d ON o.order_uid = d.order_uid
        JOIN payment p ON o.order_uid = p.order_uid
        JOIN items i ON o.order_uid = i.order_uid
        WHERE o.order_uid = $1
    `, id)
	if err != nil {
		return nil, err
	}
	return &o, nil
}
