package service

import (
	"context"

	"itk/internal/repo"
	"itk/models"

	"github.com/google/uuid"
)

type WalletService interface {
	ProcessOperation(ctx context.Context, req models.WalletOperationRequest) (int64, error)
	GetBalance(ctx context.Context, id uuid.UUID) (int64, error)
}

type walletService struct {
	repo repo.WalletRepository
}

func NewWalletService(r repo.WalletRepository) WalletService {
	return &walletService{repo: r}
}
