package repository

import (
	"back/pkg/order"

	"github.com/jmoiron/sqlx"
)

type GetOrderById interface {
	GetById(id string) (*order.OrderData, error)
}

type Repository struct {
	GetOrderById
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		GetOrderById: NewGetOrderByIdPostgres(db),
	}
}
