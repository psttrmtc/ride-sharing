package service

import (
	"context"
	"fmt"
	"ride-sharing/services/payment-service/internal/domain"
	"ride-sharing/services/payment-service/pkg/types"
	"time"

	"github.com/google/uuid"
)

type paymentService struct {
	paymentProcessor domain.PaymentProcessor
}

// NewPaymentService creates a new payment service instance
func NewPaymentService() domain.Service {
	return &paymentService{}
}

func (s *paymentService) CreatePaymentSession(
	ctx context.Context,
	tripID string,
	userID string,
	driverID string,
	amount int64,
	currency string,
) (*types.PaymentIntent, error) {
	metadata := map[string]string{
		"trip_id":   tripID,
		"user_id":   userID,
		"driver_id": driverID,
	}

	sessionID, err := s.paymentProcessor.CreatePaymentSession(ctx, amount, currency, metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment session: %w", err)
	}

	paymentIntent := &types.PaymentIntent{
		ID:              uuid.New().String(),
		TripID:          tripID,
		UserID:          userID,
		DriverID:        driverID,
		Amount:          amount,
		Currency:        currency,
		StripeSessionID: sessionID,
		CreatedAt:       time.Now(),
	}
	return paymentIntent, nil
}
