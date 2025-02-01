package service

import (
	"fmt"

	adapter "github.com//go-banking/adapter/repository"
	"github.com/rismapa/go-banking/domain"
)

type TransactionService interface {
	CreateTransaction(trx domain.Transaction, amount float64) (*domain.Transaction, error)
	GetTransactionByAccountID(accountID string) ([]domain.Transaction, error)
	GetAllTransaction() ([]domain.Transaction, error)
}

type TransactionServiceAdapter struct {
	repo adapter.TransactionRepository
	acc  adapter.AccountRepository
}

func NewTransactionService(repo adapter.TransactionRepository, acc adapter.AccountRepository) *TransactionServiceAdapter {
	return &TransactionServiceAdapter{repo: repo, acc: acc}
}

func (s *TransactionServiceAdapter) CreateTransaction(trx domain.Transaction, amount float64) (*domain.Transaction, error) {
	/*
	 * Lakukan validasi terhadap business logic terlebih dahulu
	 * sebelum melakukan membuat / memproses transaksi
	 */
	account, err := s.acc.GetAccountByID(trx.AccountID)
	if err != nil {
		return nil, err
	}

	// receiver / destination account
	if trx.Type == "transfer" {
		_, err := s.acc.GetAccountByID(trx.DestinationAccountID)
		if err != nil {
			return nil, fmt.Errorf("destination account not found")
		}

		if account.Balance < amount {
			return nil, fmt.Errorf("insufficient balance")
		}

		if amount < 10000 {
			return nil, fmt.Errorf("minimum transfer amount is 10.000")
		}
	}

	if trx.Type == "debit" && account.Balance < trx.Nominal {
		return nil, fmt.Errorf("insufficient balance")
	}

	/*
	 * Proses membuat transaksi dan update balance dilakukan
	 * ketika semua kondisi business logic terpenuhi
	 */
	data, err := s.repo.CreateTransaction(trx, amount)
	if err != nil {
		return nil, err
	}

	err = s.repo.UpdateAccountBalance(trx.AccountID, trx.DestinationAccountID, trx.Nominal, trx.Type)
	if err != nil {
		return nil, fmt.Errorf("failed to update account balance: %v", err)
	}

	return data, nil
}

func (s *TransactionServiceAdapter) GetTransactionByAccountID(accountID string) ([]domain.Transaction, error) {
	return s.repo.GetTransactionByAccountID(accountID)
}

func (s *TransactionServiceAdapter) GetAllTransaction() ([]domain.Transaction, error) {
	return s.repo.GetAllTransaction()
}
