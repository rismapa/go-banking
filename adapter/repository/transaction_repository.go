package adapter

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/rismapa/go-banking/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TransactionRepository interface {
	CreateTransaction(trx domain.Transaction, amount float64) (*domain.Transaction, error)
	UpdateAccountBalance(accountID string, destinationAccountID string, amount float64, trxType string) error
	BeginTransaction() (*sql.Tx, error)
	GetTransactionByAccountID(accountID string) ([]domain.Transaction, error)
	GetAllTransaction() ([]domain.Transaction, error)
}

type TransactionRepositoryDB struct {
	DB *sqlx.DB
}

func NewTransactionRepositoryDB(db *sqlx.DB) *TransactionRepositoryDB {
	return &TransactionRepositoryDB{DB: db}
}

/*
 * untuk debit sebagai penarikan (sisi banking)
 */
func (t *TransactionRepositoryDB) CreateTransaction(trx domain.Transaction, amount float64) (*domain.Transaction, error) {
	trx.ID = uuid.New().String()
	trx.DateTransaction = time.Now().Format("2006-01-02 15:04:05")
	query := "INSERT INTO transaction (id, date_transaction, type, nominal, account_id, destination_account_id,note) VALUES (?,?, ?, ?, ?, ?, ?)"
	_, err := t.DB.Exec(query, trx.ID, trx.DateTransaction, trx.Type, trx.Nominal, trx.AccountID, trx.DestinationAccountID, trx.Note)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %v", err)
	}

	return &trx, nil
}

func (t *TransactionRepositoryDB) UpdateAccountBalance(accountID string, destinationAccountID string, amount float64, trxType string) error {
	var query string

	switch trxType {
	case "debit":
		query = "UPDATE accounts SET balance = balance - ? WHERE id = ?"
		_, err := t.DB.Exec(query, amount, accountID)
		if err != nil {
			return fmt.Errorf("failed to update account balance: %v", err)
		}
	case "credit":
		query = "UPDATE accounts SET balance = balance + ? WHERE id = ?"
		_, err := t.DB.Exec(query, amount, accountID)
		if err != nil {
			return fmt.Errorf("failed to update account balance: %v", err)
		}
	case "transfer":
		// Memanggil UpdateSenderAndReceiverBalance dengan receiver t
		err := t.UpdateSenderAndReceiverBalance(accountID, destinationAccountID, amount)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown transaction type: %s", trxType)
	}

	return nil
}

func (t *TransactionRepositoryDB) GetTransactionByAccountID(accountID string) ([]domain.Transaction, error) {
	var transaction []domain.Transaction

	query := "SELECT id, date_transaction, type, nominal, account_id, note FROM transaction WHERE account_id = ?"
	err := t.DB.Select(&transaction, query, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction by account id: %w", err)
	}

	if len(transaction) == 0 {
		return nil, fmt.Errorf("no transaction found for account id: %s", accountID)
	}

	return transaction, nil
}

func (t *TransactionRepositoryDB) GetAllTransaction() ([]domain.Transaction, error) {
	var transaction []domain.Transaction

	query := "SELECT id, date_transaction, type, nominal, account_id, destination_account_id, note FROM transaction"
	err := t.DB.Select(&transaction, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all transaction: %w", err)
	}

	if len(transaction) == 0 {
		return nil, fmt.Errorf("no transaction found")
	}

	return transaction, nil
}

func (r *TransactionRepositoryDB) BeginTransaction() (*sql.Tx, error) {
	// Start a new transaction
	tx, err := r.DB.Begin()
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (t *TransactionRepositoryDB) UpdateSenderAndReceiverBalance(senderID string, receiverID string, amount float64) error {
	// Update sender balance
	query := "UPDATE accounts SET balance = balance - ? WHERE id = ?"
	_, err := t.DB.Exec(query, amount, senderID)
	if err != nil {
		return err
	}

	// Update receiver balance
	query = "UPDATE accounts SET balance = balance + ? WHERE id = ?"
	_, err = t.DB.Exec(query, amount, receiverID)
	if err != nil {
		return err
	}

	return nil
}
