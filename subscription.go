package gosubscription

import (
	"time"

	"github.com/google/uuid"
)

type MonthYear string

const MonthYearFormat = "01-2006"

type Subscription struct {
	Id          int       `json:"-"`
	UserId      uuid.UUID `json:"user_id"`
	ServiceName string    `json:"service_name"`
	Price       int64     `json:"price"`
	StartDate   MonthYear `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}
