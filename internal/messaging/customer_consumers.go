package messaging

import (
	"context"
	"encoding/json"
	"fmt"
)

type OrderCreatedHandler func(ctx context.Context, orderID string, customerID string) error
type PaymentCompletedHandler func(ctx context.Context, orderID string) error

type CustomerMessaging struct {
	rabbitMQ *RabbitMQ
}

func NewCustomerMessaging(rabbitMQ *RabbitMQ) *CustomerMessaging {
	return &CustomerMessaging{
		rabbitMQ: rabbitMQ,
	}
}

func (m *CustomerMessaging) SetupOrderCreatedConsumer(handler OrderCreatedHandler) error {
	queue, err := m.rabbitMQ.DeclareQueue("customer.order.created")
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	if err := m.rabbitMQ.BindQueue(queue.Name, "order.created"); err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	return m.rabbitMQ.Consume(queue.Name, func(ctx context.Context, msg []byte) error {
		var order struct {
			ID         string `json:"id"`
			CustomerID string `json:"customer_id"`
		}

		if err := json.Unmarshal(msg, &order); err != nil {
			return fmt.Errorf("failed to unmarshal order: %w", err)
		}

		return handler(ctx, order.ID, order.CustomerID)
	})
}

func (m *CustomerMessaging) SetupPaymentCompletedConsumer(handler PaymentCompletedHandler) error {
	queue, err := m.rabbitMQ.DeclareQueue("customer.payment.completed")
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	if err := m.rabbitMQ.BindQueue(queue.Name, "payment.completed"); err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	return m.rabbitMQ.Consume(queue.Name, func(ctx context.Context, msg []byte) error {
		var payment struct {
			OrderID string `json:"order_id"`
		}

		if err := json.Unmarshal(msg, &payment); err != nil {
			return fmt.Errorf("failed to unmarshal payment: %w", err)
		}

		return handler(ctx, payment.OrderID)
	})
}
