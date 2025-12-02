package repo

import (
	"context"
	"errors"
	errs "itk/internal/errors"

	"github.com/google/uuid"
)

func (r *PostgresWalletRepository) ChangeBalance(ctx context.Context, id uuid.UUID, delta int64) (int64, error) {
	const query = `
UPDATE wallets
SET balance = balance + $2
WHERE id = $1 AND balance + $2 >= 0
RETURNING balance;
`
	var balance int64
	err := r.pool.QueryRow(ctx, query, id, delta).Scan(&balance)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return 0, err
		}
		if err.Error() == "no rows in result set" {
			if delta < 0 {
				return 0, errs.ErrInsufficientFunds
			}
			return 0, errs.ErrWalletNotFound
		}
		return 0, err
	}
	return balance, nil
}

func (r *PostgresWalletRepository) GetBalance(ctx context.Context, id uuid.UUID) (int64, error) {
	const query = `SELECT balance FROM wallets WHERE id = $1`
	var balance int64
	err := r.pool.QueryRow(ctx, query, id).Scan(&balance)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return 0, errs.ErrWalletNotFound
		}
		return 0, err
	}
	return balance, nil
}

func (r *PostgresWalletRepository) CreateWallet(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		"INSERT INTO wallets (id, balance) VALUES ($1, 0)",
		id,
	)
	return err
}
