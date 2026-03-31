package service

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	gosubscription "github.com/grancc/go-effective-mobile-test"
	"github.com/grancc/go-effective-mobile-test/pkg/repository"
)

func normalizeStartDate(s string) (string, error) {
	if s == "" {
		return "", fmt.Errorf("start_date is empty")
	}
	t, err := time.Parse("2006-01", s)
	if err != nil {
		return "", fmt.Errorf("start_date: ожидается YYYY-MM: %w", err)
	}
	return t.Format("2006-01-02"), nil
}

type SubscriptionService struct {
	repo repository.Subscriptions
}

func NewSubscriptionService(repo repository.Subscriptions) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func (s *SubscriptionService) CreateSubscription(subscription gosubscription.Subscription) (int, error) {
	normalized, err := normalizeStartDate(subscription.StartDate)
	if err != nil {
		return 0, err
	}
	subscription.StartDate = normalized
	return s.repo.CreateSubscription(subscription)
}

func (s *SubscriptionService) GetSubscriptionById(id int) (gosubscription.Subscription, error) {
	return s.repo.GetSubscriptionById(id)
}

func (s *SubscriptionService) ListSubscriptionsByUserId(userID uuid.UUID) ([]gosubscription.Subscription, error) {
	return s.repo.ListSubscriptionsByUserId(userID)
}

func (s *SubscriptionService) UpdateSubscriptionById(id int, subscription gosubscription.Subscription) (gosubscription.Subscription, error) {
	normalized, err := normalizeStartDate(subscription.StartDate)
	if err != nil {
		return gosubscription.Subscription{}, err
	}
	subscription.StartDate = normalized
	return s.repo.UpdateSubscriptionById(id, subscription)
}

func (s *SubscriptionService) DeleteSubscriptionById(id int) error {
	return s.repo.DeleteSubscriptionById(id)
}

func (s *SubscriptionService) SumSubscriptions(serviceName *string, userID *uuid.UUID, periodStart, periodEnd time.Time) (int64, error) {
	return s.repo.SumSubscriptions(serviceName, userID, periodStart, periodEnd)
}
