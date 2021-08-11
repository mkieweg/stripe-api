package payment

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type Store interface {
	AddPayment(*Payment)
	GetPayment(uuid.UUID) (*Payment, error)
	ChangeStatus()
}

type StoreImplementation struct {
	Payments map[uuid.UUID]*Payment
	Mu       sync.RWMutex
}

func (s *StoreImplementation) AddPayment(payment *Payment) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.Payments[payment.UUID] = payment
}

func (s *StoreImplementation) GetPayment(id uuid.UUID) (*Payment, error) {
	p, ok := s.Payments[id]
	if !ok {
		return nil, fmt.Errorf("payment %v not found", id)
	}
	return p, nil
}

func (s *StoreImplementation) ChangeStatus() {

}
