package repository

import (
	"github.com/jmoiron/sqlx"
)

type Subscriptions interface {
	//CreateSubscription(subscription gosubscription.Subscription) (int, error)
}

type Repository struct {
	//Subscriptions
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{}
}
