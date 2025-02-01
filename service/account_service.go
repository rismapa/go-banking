package service

import (
	"fmt"

	adapter "github.com/rismapa/go-banking/adapter/repository"
	"github.com/rismapa/go-banking/domain"
	"github.com/rismapa/go-banking/dto"
)

type AccountService interface {
	GetAccounts() ([]dto.AccountWithCustomer, error)
	CreateAccount(account domain.Account) (*domain.Account, error)
	GetAccountByID(id string) (*domain.Account, error)
	GetAccountByUsername(username string) (*domain.Account, error)
	GetAccountByCustomerID(id string) (*domain.Customer, []domain.Account, error)
	UpdateAccount(id string, account domain.Account) (*domain.Account, error)
	SoftDeleteAccount(id string) (*domain.Account, error)
}

type AccountAdapter struct {
	repo adapter.AccountRepository
	cust adapter.CustomerRepository
}

func NewAccountService(repo adapter.AccountRepository, cust adapter.CustomerRepository) *AccountAdapter {
	return &AccountAdapter{repo: repo, cust: cust}
}

func (s *AccountAdapter) CreateAccount(account domain.Account) (*domain.Account, error) {
	_, err := s.cust.GetCustomerByID(account.Customer_ID)
	if err != nil {
		return nil, err
	}

	if account.Balance < 0 {
		return nil, fmt.Errorf("account balance cannot be negative")
	}

	return s.repo.CreateAccount(account)
}

func (s *AccountAdapter) GetAccounts() ([]dto.AccountWithCustomer, error) {
	return s.repo.GetAccounts()
}

func (s *AccountAdapter) GetAccountByID(id string) (*domain.Account, error) {
	return s.repo.GetAccountByID(id)
}

func (s *AccountAdapter) GetAccountByUsername(username string) (*domain.Account, error) {
	return s.repo.GetAccountByUsername(username)
}

func (s *AccountAdapter) GetAccountByCustomerID(id string) (*domain.Customer, []domain.Account, error) {
	customer, err := s.cust.GetCustomerByID(id)
	if err != nil {
		return customer, nil, err
	}

	accounts, err := s.repo.GetAccountByCustomerID(id)
	if err != nil {
		return nil, accounts, err
	}

	return customer, accounts, nil
}

func (s *AccountAdapter) UpdateAccount(id string, account domain.Account) (*domain.Account, error) {
	_, err := s.repo.GetAccountByID(id)
	if err != nil {
		return nil, err
	}

	_, err = s.cust.GetCustomerByID(account.Customer_ID)
	if err != nil {
		return nil, err
	}

	account.ID = id
	updatedData, err := s.repo.UpdateAccount(account)
	if err != nil {
		return nil, fmt.Errorf("internal error while updating customer: %v", err)
	}

	return updatedData, nil
}

func (s *AccountAdapter) SoftDeleteAccount(id string) (*domain.Account, error) {
	existAccount, err := s.repo.GetAccountByID(id)
	if err != nil {
		return nil, fmt.Errorf("internal error while checking account: %v", err)
	}

	if !existAccount.Status {
		return nil, fmt.Errorf("account already deleted")
	}

	return s.repo.SoftDeleteAccount(*existAccount)
}
