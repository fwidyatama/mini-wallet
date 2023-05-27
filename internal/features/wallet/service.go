package wallet

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"mini-wallet/domain/constant"
	historyModel "mini-wallet/domain/model/history"
	model "mini-wallet/domain/model/wallet"
)

type walletService struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &walletService{repository}
}

type Service interface {
	EnableWallet(ctx context.Context, args model.OwnerWalletParam) (*model.Wallet, error)
	DisableWallet(ctx context.Context, args model.OwnerWalletParam) (*model.Wallet, error)
	GetWallet(ctx context.Context, args model.OwnerWalletParam) (*model.Wallet, error)
	DepositWallet(ctx context.Context, args historyModel.TransactionParams) (*historyModel.Deposit, error)
	WithdrawalWallet(ctx context.Context, args historyModel.TransactionParams) (*historyModel.Withdrawal, error)
	GetTransaction(ctx context.Context, args historyModel.TransactionParams) (*[]historyModel.History, error)
}

func (s *walletService) EnableWallet(ctx context.Context, args model.OwnerWalletParam) (*model.Wallet, error) {
	existingWallet, _ := s.repository.GetWalletByOwner(ctx, args.Owner)

	// if owner doesn't have wallet, create it
	if existingWallet.OwnedBy == "" {
		result, err := s.repository.CreateWallet(ctx, args)
		if err != nil {
			logrus.Errorf("Service error : %v", err)
			return nil, err
		}
		return result, nil
	} else {
		// if status already enabled, return error
		if existingWallet.Status == constant.Enabled {
			err := errors.New(constant.AlreadyEnabled)
			return nil, err
		} else {
			// else if status disabled, enable it.
			result, err := s.repository.EnabledWallet(ctx, args)
			if err != nil {
				logrus.Errorf("Service error : %v", err)
				return nil, err
			}
			return result, nil
		}
	}
}

func (s *walletService) DisableWallet(ctx context.Context, args model.OwnerWalletParam) (*model.Wallet, error) {
	result, err := s.repository.DisabledWallet(ctx, args)
	if err != nil {
		logrus.Errorf("Service error : %v", err)
		return nil, err
	}

	return result, nil
}

func (s *walletService) GetWallet(ctx context.Context, args model.OwnerWalletParam) (*model.Wallet, error) {
	wallet, err := s.repository.GetWalletByOwner(ctx, args.Owner)
	if err != nil {
		logrus.Errorf("Service error : %v", err)
		return nil, err
	}

	if wallet.Status == constant.Disabled {
		err = errors.New(constant.WalletDisabled)
		logrus.Errorf("service error : %v", err)
	}

	return wallet, nil
}
func (s *walletService) DepositWallet(ctx context.Context, args historyModel.TransactionParams) (*historyModel.Deposit, error) {
	transaction, err := s.repository.DepositWallet(ctx, args)

	if err != nil {
		logrus.Errorf("Service error : %v", err)
		return nil, err
	}

	return transaction, nil
}

func (s *walletService) WithdrawalWallet(ctx context.Context, args historyModel.TransactionParams) (*historyModel.Withdrawal, error) {
	transaction, err := s.repository.WithdrawalWallet(ctx, args)

	if err != nil {
		logrus.Errorf("Service error : %v", err)
		return nil, err
	}

	return transaction, nil
}

func (s *walletService) GetTransaction(ctx context.Context, args historyModel.TransactionParams) (*[]historyModel.History, error) {
	wallet, _ := s.GetWallet(ctx, model.OwnerWalletParam{Owner: args.TransactionBy})

	if wallet.Status == constant.Disabled {
		err := errors.New(constant.WalletDisabled)
		logrus.Errorf("service error : %v", err)
		return nil, err
	} else {
		transaction, err := s.repository.GetTransaction(ctx, args)

		if err != nil {
			logrus.Errorf("Service error : %v", err)
			return nil, err
		}

		return transaction, nil
	}

}
