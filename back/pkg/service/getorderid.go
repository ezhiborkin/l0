package service

import (
	"back/pkg/order"
	"back/pkg/repository"
	"errors"
	"fmt"
)

type GetOrderPService struct {
	repo repository.GetOrderById
}

func NewGetOrderPService(repo repository.GetOrderById) *GetOrderPService {
	return &GetOrderPService{
		repo: repo,
	}
}

func (s *GetOrderPService) GetById(id string) (*order.OrderData, error) {
	order, err := s.repo.GetById(id)
	if err != nil {
		fmt.Println("getorderid")
		return nil, errors.New("geto")
	}

	return order, nil
}
