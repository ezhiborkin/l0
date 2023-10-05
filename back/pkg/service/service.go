package service

import (
	"back/pkg/cashe"
	"back/pkg/order"
	"back/pkg/repository"
)

type GetOrder interface {
	GetById(id string) (*order.OrderData, error)
}

type Service struct {
	GetOrder
}

func NewService(repos *repository.Repository, cache *cashe.Cache) *Service {
	return &Service{}
}
