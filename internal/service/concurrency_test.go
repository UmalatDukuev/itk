package service

import (
	"context"
	"sync"
	"testing"
	"time"

	"itk/models"

	"github.com/google/uuid"
)

func TestWalletService_Concurrent1000(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r := newMockRepo()
	s := NewWalletService(r)

	id := uuid.New()
	if err := r.CreateWallet(ctx, id); err != nil {
		t.Fatalf("failed to create wallet: %v", err)
	}

	const workers = 1000
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func(i int) {
			defer wg.Done()

			req := models.WalletOperationRequest{
				WalletID:      id,
				Amount:        1,
				OperationType: models.OperationDeposit,
			}

			_, _ = s.ProcessOperation(ctx, req)
		}(i)
	}

	wg.Wait()

	bal, err := s.GetBalance(ctx, id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if bal != workers {
		t.Fatalf("expected balance %d, got %d", workers, bal)
	}
}
