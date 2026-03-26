package repository

import (
	"time"

	"github.com/google/uuid"
	gosubscription "github.com/grancc/go-effective-mobile-test"
	"github.com/jmoiron/sqlx"
)

type Subscriptions interface {
	CreateSubscription(subscription gosubscription.Subscription) (int, error)
	GetSubscriptionById(id int) (gosubscription.Subscription, error)
	ListSubscriptionsByUserId(userID uuid.UUID) ([]gosubscription.Subscription, error)
	UpdateSubscriptionById(id int, subscription gosubscription.Subscription) (gosubscription.Subscription, error)
	DeleteSubscriptionById(id int) error
	SumSubscriptions(serviceName *string, userID *uuid.UUID, periodStart, periodEnd time.Time) (int64, error)
}

type Repository struct {
	Subscriptions
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Subscriptions: NewSubsPostgres(db),
	}
}
