package service

import (
	"context"
	"errors"
	errs "itk/internal/errors"
	"itk/models"
	"sync"
	"testing"

	"github.com/google/uuid"
)

type mockRepo struct {
	mu        sync.Mutex
	balances  map[uuid.UUID]int64
	forceErr  error
	forceGet  error
	forceChan error
}

func newMockRepo() *mockRepo {
	return &mockRepo{balances: make(map[uuid.UUID]int64)}
}

func (m *mockRepo) CreateWallet(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.balances[id] = 0
	return nil
}

func (m *mockRepo) ChangeBalance(ctx context.Context, id uuid.UUID, delta int64) (int64, error) {
	if m.forceChan != nil {
		return 0, m.forceChan
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	bal, ok := m.balances[id]
	if !ok {
		return 0, errs.ErrWalletNotFound
	}

	if bal+delta < 0 {
		return 0, errs.ErrInsufficientFunds
	}

	bal += delta
	m.balances[id] = bal
	return bal, nil
}

func (m *mockRepo) GetBalance(ctx context.Context, id uuid.UUID) (int64, error) {
	if m.forceGet != nil {
		return 0, m.forceGet
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	bal, ok := m.balances[id]
	if !ok {
		return 0, errs.ErrWalletNotFound
	}
	return bal, nil
}

func TestProcessOperation_Deposit(t *testing.T) {
	r := newMockRepo()
	id := uuid.New()
	r.balances[id] = 0

	s := NewWalletService(r)

	req := models.WalletOperationRequest{
		WalletID:      id,
		OperationType: models.OperationDeposit,
		Amount:        100,
	}

	bal, err := s.ProcessOperation(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if bal != 100 {
		t.Fatalf("expected 100, got %d", bal)
	}
}

func TestProcessOperation_Withdraw(t *testing.T) {
	r := newMockRepo()
	id := uuid.New()
	r.balances[id] = 200

	s := NewWalletService(r)

	req := models.WalletOperationRequest{
		WalletID:      id,
		OperationType: models.OperationWithdraw,
		Amount:        50,
	}

	bal, err := s.ProcessOperation(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if bal != 150 {
		t.Fatalf("expected 150, got %d", bal)
	}
}

func TestProcessOperation_InvalidAmount(t *testing.T) {
	r := newMockRepo()
	s := NewWalletService(r)

	req := models.WalletOperationRequest{
		WalletID:      uuid.New(),
		OperationType: models.OperationDeposit,
		Amount:        0,
	}

	_, err := s.ProcessOperation(context.Background(), req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestProcessOperation_InvalidOperationType(t *testing.T) {
	r := newMockRepo()
	s := NewWalletService(r)

	req := models.WalletOperationRequest{
		WalletID:      uuid.New(),
		OperationType: "WHAT_IS_THIS",
		Amount:        10,
	}

	_, err := s.ProcessOperation(context.Background(), req)
	if !errors.Is(err, errs.ErrInvalidOperationType) {
		t.Fatalf("wrong error: %v", err)
	}
}

func TestProcessOperation_RepoError(t *testing.T) {
	r := newMockRepo()
	id := uuid.New()
	r.balances[id] = 0
	r.forceChan = errors.New("boom")

	s := NewWalletService(r)

	req := models.WalletOperationRequest{
		WalletID:      id,
		OperationType: models.OperationDeposit,
		Amount:        10,
	}

	_, err := s.ProcessOperation(context.Background(), req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetBalance_OK(t *testing.T) {
	r := newMockRepo()
	id := uuid.New()
	r.balances[id] = 123

	s := NewWalletService(r)

	bal, err := s.GetBalance(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if bal != 123 {
		t.Fatalf("expected 123, got %d", bal)
	}
}

func TestGetBalance_Error(t *testing.T) {
	r := newMockRepo()
	id := uuid.New()
	r.forceGet = errors.New("db down")

	s := NewWalletService(r)

	_, err := s.GetBalance(context.Background(), id)
	if err == nil {
		t.Fatal("expected error")
	}
}
