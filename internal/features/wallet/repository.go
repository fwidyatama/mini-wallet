package wallet

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"mini-wallet/domain/constant"
	historyModel "mini-wallet/domain/model/history"
	"mini-wallet/domain/model/wallet"
)

type walletRepository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &walletRepository{db}
}

type Repository interface {
	CreateWallet(ctx context.Context, args model.OwnerWalletParam) (*model.Wallet, error)
	GetWalletByOwner(ctx context.Context, ownerId string) (*model.Wallet, error)
	DisabledWallet(ctx context.Context, args model.OwnerWalletParam) (*model.Wallet, error)
	EnabledWallet(ctx context.Context, args model.OwnerWalletParam) (*model.Wallet, error)

	DepositWallet(ctx context.Context, args historyModel.TransactionParams) (*historyModel.Deposit, error)
	WithdrawalWallet(ctx context.Context, args historyModel.TransactionParams) (*historyModel.Withdrawal, error)
	UpdateBalance(ctx context.Context, args model.UpdateBalanceParam) error
	GetTransaction(ctx context.Context, args historyModel.TransactionParams) (*[]historyModel.History, error)
}

func (r *walletRepository) CreateWallet(ctx context.Context, args model.OwnerWalletParam) (*model.Wallet, error) {
	var result model.Wallet
	err := r.db.QueryRowxContext(ctx, InsertWalletQuery, args.Owner).Scan(
		&result.ID,
		&result.OwnedBy,
		&result.EnabledAt,
		&result.Balance,
		&result.Status,
	)
	if err != nil {
		logrus.Errorf("Repository error : %v", err.Error())
		return nil, err
	}

	return &result, nil
}

func (r *walletRepository) GetWalletByOwner(ctx context.Context, ownerId string) (*model.Wallet, error) {
	var result model.Wallet

	err := r.db.QueryRowxContext(ctx, GetWalletByOwnerIdQuery, ownerId).Scan(
		&result.ID,
		&result.OwnedBy,
		&result.Status,
		&result.Balance,
		&result.EnabledAt,
		&result.DisabledAt,
	)
	if err != nil && err != sql.ErrNoRows {
		logrus.Errorf("Repository error : %v", err.Error())
		return nil, err
	}

	return &result, nil
}

func (r *walletRepository) DisabledWallet(ctx context.Context, args model.OwnerWalletParam) (*model.Wallet, error) {
	var result model.Wallet
	err := r.db.QueryRowxContext(ctx, DisableWalletQuery, args.Owner).Scan(
		&result.ID,
		&result.OwnedBy,
		&result.DisabledAt,
		&result.Balance,
		&result.Status,
	)
	if err != nil {
		logrus.Errorf("Repository error : %v", err.Error())
		return nil, err
	}

	return &result, nil
}

func (r *walletRepository) EnabledWallet(ctx context.Context, args model.OwnerWalletParam) (*model.Wallet, error) {
	var result model.Wallet
	err := r.db.QueryRowxContext(ctx, EnableWalletQuery, args.Owner).Scan(
		&result.ID,
		&result.OwnedBy,
		&result.DisabledAt,
		&result.Balance,
		&result.Status,
	)
	if err != nil {
		logrus.Errorf("Repository error : %v", err.Error())
		return nil, err
	}

	return &result, nil
}

func (r *walletRepository) UpdateBalance(ctx context.Context, args model.UpdateBalanceParam) error {
	_, err := r.db.ExecContext(ctx, UpdateBalanceWalletQuery, args.Balance, args.WalletId)
	if err != nil {
		logrus.Errorf("Repository error : %v", err.Error())
		return err
	}
	return nil
}

func (r *walletRepository) DepositWallet(ctx context.Context, args historyModel.TransactionParams) (*historyModel.Deposit, error) {
	var result historyModel.Deposit
	var wallet model.Wallet

	ownerId := ctx.Value("ownerId").(string)

	// begin transaction
	tx, err := r.db.Begin()

	// get wallet by owner id
	err = tx.QueryRowContext(ctx, GetWalletByOwnerIdQuery, ownerId).Scan(
		&wallet.ID,
		&wallet.OwnedBy,
		&wallet.Status,
		&wallet.Balance,
		&wallet.EnabledAt,
		&wallet.DisabledAt,
	)

	if wallet.Status != constant.Disabled {
		if err != nil {
			logrus.Error(err)
		}

		// create transaction
		err = tx.QueryRowContext(ctx, TransactionWalletQuery,
			wallet.ID,
			constant.SuccessMessage,
			ownerId,
			constant.Deposit,
			args.Amount,
			args.ReferenceID).Scan(
			&result.ID,
			&result.Status,
			&result.DepositedBy,
			&result.DepositedAt,
			&result.Amount,
			&result.ReferenceId)

		if err != nil {
			logrus.Error(err)
		}

		// add balance with amount
		balance := wallet.Balance + args.Amount

		// update balance wallet
		_, err = tx.ExecContext(ctx, UpdateBalanceWalletQuery, balance, wallet.ID)

		err = tx.Commit()

		if err != nil {
			defer tx.Rollback()
			return nil, err
		}

		return &result, nil
	} else {
		err = errors.New(constant.WalletDisabled)
		return nil, err
	}
}

func (r *walletRepository) WithdrawalWallet(ctx context.Context, args historyModel.TransactionParams) (*historyModel.Withdrawal, error) {
	var result historyModel.Withdrawal
	var wallet model.Wallet

	ownerId := ctx.Value("ownerId").(string)

	// begin transaction
	tx, err := r.db.Begin()

	// get wallet by owner id
	err = tx.QueryRowContext(ctx, GetWalletByOwnerIdQuery, ownerId).Scan(
		&wallet.ID,
		&wallet.OwnedBy,
		&wallet.Status,
		&wallet.Balance,
		&wallet.EnabledAt,
		&wallet.DisabledAt,
	)
	if err != nil {
		logrus.Error(err)
	}

	if wallet.Status != constant.Disabled {

		// create transaction
		err = tx.QueryRowContext(ctx, TransactionWalletQuery,
			wallet.ID,
			constant.SuccessMessage,
			ownerId,
			constant.Withdrawal,
			args.Amount,
			args.ReferenceID).Scan(
			&result.ID,
			&result.Status,
			&result.WithdrawnBy,
			&result.WithdrawnAt,
			&result.Amount,
			&result.ReferenceId)

		if err != nil {
			logrus.Error(err)
		}

		if wallet.Balance < args.Amount {
			tx.Rollback()
			err = errors.New(constant.AmountExceeded)
			return nil, err
		}

		// add balance with amount
		balance := wallet.Balance - args.Amount

		// update balance wallet
		_, err = tx.ExecContext(ctx, UpdateBalanceWalletQuery, balance, wallet.ID)

		err = tx.Commit()

		if err != nil && err.Error() != constant.AmountExceeded {
			defer tx.Rollback()
			return nil, err
		}
		return &result, nil
	} else {
		err = errors.New(constant.WalletDisabled)
		return nil, err
	}

}

func (r *walletRepository) GetTransaction(ctx context.Context, args historyModel.TransactionParams) (*[]historyModel.History, error) {
	var results []historyModel.History

	rows, err := r.db.QueryContext(ctx, GetWalletTransactionQuery, args.TransactionBy)
	if err != nil {
		logrus.Errorf("Repository error : %v", err.Error())
		return nil, err
	}
	for rows.Next() {
		var transaction historyModel.History
		err = rows.Scan(&transaction.ID, &transaction.WalletID, &transaction.Status, &transaction.TransactionBy, &transaction.Type, &transaction.Amount, &transaction.ReferenceID, &transaction.TransactionAt, &transaction.CreatedAt)
		results = append(results, transaction)
	}

	defer rows.Close()

	return &results, nil
}
