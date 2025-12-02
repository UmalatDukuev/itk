package service

import (
	"context"
	"errors"
	errs "itk/internal/errors"
	"itk/models"

	"github.com/google/uuid"
)

func (s *walletService) ProcessOperation(ctx context.Context, req models.WalletOperationRequest) (int64, error) {
	if req.Amount <= 0 {
		return 0, errors.New("amount must be positive")
	}

	var delta int64
	switch req.OperationType {
	case models.OperationDeposit:
		delta = req.Amount
	case models.OperationWithdraw:
		delta = -req.Amount
	default:
		return 0, errs.ErrInvalidOperationType
	}

	return s.repo.ChangeBalance(ctx, req.WalletID, delta)
}

func (s *walletService) GetBalance(ctx context.Context, id uuid.UUID) (int64, error) {
	return s.repo.GetBalance(ctx, id)
}
