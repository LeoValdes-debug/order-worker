package notifier

import (
	"context"
	"log"
)

// stub - в реальном проекте тут был бы HTTP-вызов к сервису уведомлений
type Notifier struct{}

func New() *Notifier {
	return &Notifier{}
}

func (n *Notifier) Send(_ context.Context, userID int, message string) error {
	log.Printf("notify user_id=%d: %s", userID, message)
	return nil
}
