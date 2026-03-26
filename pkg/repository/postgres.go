package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	gosubscription "github.com/grancc/go-effective-mobile-test"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
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
	query := fmt.Sprintf("INSERT INTO %s (service_name, price, user_id, start_date, end_date) values ($1, $2, $3, $4, $5) RETURNING id",
		userSubsTable)
	row := a.db.QueryRow(query, &uSub.ServiceName, &uSub.Price, &uSub.UserId, &uSub.StartDate, &uSub.EndDate)
	if err := row.Scan(&id); err != nil {
		logrus.WithError(err).WithField("op", "CreateSubscription").Error("db")
		return 0, err
	}
	logrus.WithFields(logrus.Fields{"op": "CreateSubscription", "id": id}).Debug("ok")
	return id, nil
}

func (a *SubsPostgres) GetSubscriptionById(id int) (gosubscription.Subscription, error) {
	var sub gosubscription.Subscription
	query := fmt.Sprintf("select * from %s where id = $1", userSubsTable)
	row := a.db.QueryRow(query, id)

	if err := row.Scan(&sub.Id, &sub.UserId, &sub.ServiceName, &sub.Price, &sub.StartDate, &sub.EndDate); err != nil {
		if err == sql.ErrNoRows {
			return sub, err
		}
		logrus.WithError(err).WithFields(logrus.Fields{"op": "GetSubscriptionById", "id": id}).Error("db")
		return sub, err
	}
	return sub, nil
}

func (a *SubsPostgres) UpdateSubscriptionById(id int, sub gosubscription.Subscription) (gosubscription.Subscription, error) {
	var sub_rez gosubscription.Subscription
	query := `update subscription set 
				service_name = $1,
				price = $2,
				start_date = $3,
				end_date = $4
				where id = $5
				RETURNING id, user_id, service_name, price, start_date, end_date;`

	row := a.db.QueryRow(query, &sub.ServiceName, &sub.Price, &sub.StartDate, &sub.EndDate, id)

	if err := row.Scan(&sub_rez.Id, &sub_rez.UserId, &sub_rez.ServiceName, &sub_rez.Price, &sub_rez.StartDate, &sub_rez.EndDate); err != nil {
		if err == sql.ErrNoRows {
			return sub_rez, err
		}
		logrus.WithError(err).WithFields(logrus.Fields{"op": "UpdateSubscriptionById", "id": id}).Error("db")
		return sub_rez, err
	}
	return sub_rez, nil
}

func (a *SubsPostgres) DeleteSubscriptionById(id int) error {
	query := fmt.Sprintf("delete from %s where id = $1", userSubsTable)
	_, err := a.db.Exec(query, id)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{"op": "DeleteSubscriptionById", "id": id}).Error("db")
		return err
	}
	return nil
}

func (a *SubsPostgres) ListSubscriptionsByUserId(userID uuid.UUID) ([]gosubscription.Subscription, error) {

	query := fmt.Sprintf("select id, user_id, service_name, price, start_date, end_date from %s where user_id = $1", userSubsTable)
	rows, err := a.db.Query(query, userID)
	if err != nil {
		logrus.WithError(err).WithField("op", "ListSubscriptionsByUserId").Error("db")
		return nil, err
	}
	defer rows.Close()

	subs := []gosubscription.Subscription{}
	for rows.Next() {
		p := gosubscription.Subscription{}
		err := rows.Scan(&p.Id, &p.UserId, &p.ServiceName, &p.Price, &p.StartDate, &p.EndDate)
		if err != nil {
			logrus.WithError(err).WithField("op", "ListSubscriptionsByUserId").Error("scan")
			return nil, err
		}
		subs = append(subs, p)
	}
	if err := rows.Err(); err != nil {
		logrus.WithError(err).WithField("op", "ListSubscriptionsByUserId").Error("rows")
		return nil, err
	}
	return subs, nil
}

func (a *SubsPostgres) SumSubscriptions(serviceName *string, userID *uuid.UUID, periodStart, periodEnd time.Time) (int64, error) {
	query := fmt.Sprintf(`
		SELECT COALESCE(SUM(
			(price * (
				EXTRACT(YEAR FROM age(hi, lo)) * 12 +
				EXTRACT(MONTH FROM age(hi, lo))
			))::bigint
		), 0)
		FROM (
			SELECT
				price,
				GREATEST(to_date(start_date, 'MM-YYYY'), $1::date) AS lo,
				LEAST(COALESCE(end_date, 'infinity'::date), $2::date) AS hi
			FROM %s
			WHERE to_date(start_date, 'MM-YYYY') <= $2::date
			  AND COALESCE(end_date, 'infinity'::date) >= $1::date
			  AND ($3::uuid IS NULL OR user_id = $3)
			  AND ($4::text IS NULL OR service_name = $4)
		) AS t
		WHERE lo <= hi
	`, userSubsTable)

	var total int64
	err := a.db.QueryRow(query, periodStart, periodEnd, userID, serviceName).Scan(&total)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"op": "SumSubscriptions",
		}).Error("db")
		return 0, err
	}
	logrus.WithFields(logrus.Fields{
		"op":    "SumSubscriptions",
		"total": total,
	}).Debug("ok")
	return total, nil
}
