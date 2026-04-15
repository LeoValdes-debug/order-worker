package processor

import (
	"context"
	"fmt"
	"log"

	"github.com/leovaldes-debug/order-worker/internal/consumer"
	"github.com/leovaldes-debug/order-worker/internal/notifier"
)

type Processor struct {
	notifier *notifier.Notifier
}

func New(n *notifier.Notifier) *Processor {
	return &Processor{notifier: n}
}

func (p *Processor) Handle(ctx context.Context, event consumer.OrderEvent) error {
	log.Printf("processing order_id=%d status=%s", event.OrderID, event.Status)

	switch event.Status {
	case "created":
		return p.handleCreated(ctx, event)
	case "cancelled":
		return p.handleCancelled(ctx, event)
	default:
		log.Printf("unknown status %q, skipping", event.Status)
		return nil
	}
}

func (p *Processor) handleCreated(ctx context.Context, event consumer.OrderEvent) error {
	msg := fmt.Sprintf("Order #%d confirmed. Total: %.2f", event.OrderID, event.TotalPrice)
	if err := p.notifier.Send(ctx, event.UserID, msg); err != nil {
		return fmt.Errorf("notify created: %w", err)
	}
	return nil
}

func (p *Processor) handleCancelled(ctx context.Context, event consumer.OrderEvent) error {
	msg := fmt.Sprintf("Order #%d has been cancelled.", event.OrderID)
	if err := p.notifier.Send(ctx, event.UserID, msg); err != nil {
		return fmt.Errorf("notify cancelled: %w", err)
	}
	return nil
}
