package repository

import (
	"fmt"

	gosubscription "github.com/grancc/go-effective-mobile-test"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	userSubsTable = "subscription"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

type SubsPostgres struct {
	db *sqlx.DB
}

func NewSubsPostgres(db *sqlx.DB) *SubsPostgres {
	return &SubsPostgres{db: db}
}

func (a *SubsPostgres) CreateSubscription(uSub gosubscription.Subscription) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (service_name, price, user_id, start_date) values ($1, $2, $3, $4) RETURNING id",
		userSubsTable)
	row := a.db.QueryRow(query, &uSub.ServiceName, &uSub.Price, &uSub.UserId, &uSub.StartDate)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
