package service

import (
	"back/pkg/order"
	"back/pkg/repository"
)

type GetOrderService struct {
	repo repository.GetOrderById
}

func NewGetOrderService(repo repository.GetOrderById) *GetOrderService {
	return &GetOrderService{repo: repo}
}

func (s *GetOrderService) GetById(id string) (*order.OrderData, error) {
	order, err := s.repo.GetById(id)
	if err != nil {
		return nil, err
	}

	return order, nil
}
