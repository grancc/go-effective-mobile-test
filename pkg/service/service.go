package service

import (
	"time"

	"github.com/google/uuid"
	gosubscription "github.com/grancc/go-effective-mobile-test"
	"github.com/grancc/go-effective-mobile-test/pkg/repository"
)

type Subscriptions interface {
	CreateSubscription(subscription gosubscription.Subscription) (int, error)
	GetSubscriptionById(id int) (gosubscription.Subscription, error)
	ListSubscriptionsByUserId(userID uuid.UUID) ([]gosubscription.Subscription, error)
	UpdateSubscriptionById(id int, subscription gosubscription.Subscription) (gosubscription.Subscription, error)
	DeleteSubscriptionById(id int) error
	SumSubscriptions(serviceName *string, userID *uuid.UUID, periodStart, periodEnd time.Time) (int64, error)
}

type Service struct {
	Subscriptions
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Subscriptions: repo.Subscriptions,
	}
}
