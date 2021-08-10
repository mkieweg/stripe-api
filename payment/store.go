package payment

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type Store interface {
	AddPayment(Payment)
	GetPayment(uuid.UUID) (*Payment, error)
	ChangeStatus()
}

type StoreImplementation struct {
	payments map[uuid.UUID]Payment
	mu       sync.RWMutex
}

func (s *StoreImplementation) AddPayment(payment Payment) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.payments[payment.UUID] = payment
}

func (s *StoreImplementation) GetPayment(id uuid.UUID) (*Payment, error) {
	p, ok := s.payments[id]
	if !ok {
		return nil, fmt.Errorf("payment %v not found", id)
	}
	return &p, nil
}

func (s *StoreImplementation) ChangeStatus() {

}
