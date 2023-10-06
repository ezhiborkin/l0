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
	var ord order.OrderData

	query := `
    SELECT o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature, o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard, d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
    p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee,
    i.chrt_id, i.track_number, i.price, i.rid, i.name, i.sale, i.size, i.total_price, i.nm_id, i.brand, i.status
    FROM orders o
    JOIN delivery d ON o.order_uid = d.order_uid
    JOIN payment p ON o.order_uid = p.order_uid
    JOIN items i ON o.order_uid = i.order_uid
    WHERE o.order_uid = $1
    `

	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
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
	}

	return &ord, nil
}
