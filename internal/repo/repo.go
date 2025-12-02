package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WalletRepository interface {
	ChangeBalance(ctx context.Context, id uuid.UUID, delta int64) (int64, error)
	GetBalance(ctx context.Context, id uuid.UUID) (int64, error)
	CreateWallet(ctx context.Context, id uuid.UUID) error
}

type PostgresWalletRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresWalletRepository(pool *pgxpool.Pool) *PostgresWalletRepository {
	return &PostgresWalletRepository{pool: pool}
}
