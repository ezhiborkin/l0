package service

import (
	"back/pkg/order"
	"back/pkg/repository"
)

type GetOrderP interface {
	GetById(id string) (*order.OrderData, error)
}

type Service struct {
	GetOrderP
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		GetOrderP: NewGetOrderPService(repos.GetOrderById),
	}
}
