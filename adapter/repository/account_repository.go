package adapter

import (
	"fmt"

	"github.com/rismapa/go-banking/domain"
	"github.com/rismapa/go-banking/dto"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AccountRepository interface {
	CreateAccount(account domain.Account) (*domain.Account, error)
	GetAccounts() ([]dto.AccountWithCustomer, error)
	GetAccountByID(id string) (*domain.Account, error)
	GetAccountByUsername(username string) (*domain.Account, error)
	GetAccountByCustomerID(id string) ([]domain.Account, error)
	UpdateAccount(account domain.Account) (*domain.Account, error)
	SoftDeleteAccount(account domain.Account) (*domain.Account, error)
}

type AccountRepositoryDB struct {
	DB *sqlx.DB
}

func NewAccountRepositoryDB(db *sqlx.DB) *AccountRepositoryDB {
	return &AccountRepositoryDB{DB: db}
}

func (a *AccountRepositoryDB) CreateAccount(account domain.Account) (*domain.Account, error) {
	account.ID = uuid.New().String()

	query := "INSERT INTO accounts (id, customer_id, username, password, balance, currency, status) VALUES (?, ?, ?, ?, ?, ?, ?)"
	_, err := a.DB.Exec(query, account.ID, account.Customer_ID, account.Username, account.Password, account.Balance, account.Currency, account.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %v", err)
	}

	return &account, nil
}

func (a *AccountRepositoryDB) GetAccounts() ([]dto.AccountWithCustomer, error) {
	var accounts []dto.AccountWithCustomer
	query := "SELECT accounts.id, customers.name, accounts.username, accounts.balance, accounts.currency, accounts.status FROM accounts INNER JOIN customers ON accounts.customer_id = customers.id"

	err := a.DB.Select(&accounts, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts: %v", err)
	}

	if len(accounts) == 0 {
		return nil, fmt.Errorf("no accounts found")
	}

	return accounts, nil
}

func (a *AccountRepositoryDB) GetAccountByID(id string) (*domain.Account, error) {
	var account domain.Account
	query := "SELECT id, customer_id, username, balance, currency, status FROM accounts WHERE id = ?"
	err := a.DB.Get(&account, query, id)
	if err != nil {
		if account == (domain.Account{}) {
			return nil, fmt.Errorf("no accounts found")
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	return &account, nil
}

func (a *AccountRepositoryDB) GetAccountByUsername(username string) (*domain.Account, error) {
	var account domain.Account
	query := "SELECT id, customer_id, username, password, balance, currency, status FROM accounts WHERE username = ?"
	err := a.DB.Get(&account, query, username)
	if err != nil {
		if account == (domain.Account{}) {
			return nil, fmt.Errorf("no accounts found")
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	return &account, nil
}

func (a *AccountRepositoryDB) GetAccountByCustomerID(id string) ([]domain.Account, error) {
	accountQuery := "SELECT id, customer_id, username, balance, currency, status FROM accounts WHERE customer_id = ?"

	var accounts []domain.Account
	err := a.DB.Select(&accounts, accountQuery, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get account base on customer id: %v", err)
	}

	if len(accounts) == 0 {
		return nil, fmt.Errorf("no accounts found based on customer id: %v", err)
	}

	return accounts, nil
}

func (a *AccountRepositoryDB) UpdateAccount(account domain.Account) (*domain.Account, error) {
	query := "UPDATE accounts SET username = ?, password = ?, balance = ?, currency = ?, status = ? WHERE id = ?"
	_, err := a.DB.Exec(query, account.Username, account.Password, account.Balance, account.Currency, account.Status, account.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update account: %v", err)
	}

	return &account, nil
}

func (a *AccountRepositoryDB) SoftDeleteAccount(account domain.Account) (*domain.Account, error) {
	query := "UPDATE accounts SET status = 0 WHERE id = ?"
	_, err := a.DB.Exec(query, account.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete account: %v", err)
	}

	account.Status = false

	return &account, nil
}
