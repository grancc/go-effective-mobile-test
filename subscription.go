package gosubscription

import (
	"github.com/google/uuid"
)

// Subscription модель подписки (тело запроса/ответа в API).
type Subscription struct {
	Id          int       `json:"id" db:"id"`
	UserId      uuid.UUID `json:"user_id" db:"user_id"`
	ServiceName string    `json:"service_name" db:"service_name"`
	Price       int64     `json:"price" db:"price"`
	StartDate   string    `json:"start_date" db:"start_date"`
	EndDate     string    `json:"end_date" db:"end_date"`
}
