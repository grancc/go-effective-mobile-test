package service

import (
	"github.com/grancc/go-effective-mobile-test/pkg/repository"
)

type Subscriptions interface {
	//CreateSubscription(subscription gosubscription.Subscription) (int, error)
}

type Service struct {
	//Subscriptions
}

func NewService(repo *repository.Repository) *Service {
	return &Service{}
}
