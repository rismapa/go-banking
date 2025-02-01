package adapter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rismapa/go-banking/domain"
	"github.com/rismapa/go-banking/dto"
	"github.com/rismapa/go-banking/service"
	"github.com/rismapa/go-banking/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	logger "github.com/rismapa/go-banking-lib/config"
)

type TransactionHandlerDB struct {
	Service   service.TransactionService
	Validator validator.Validate
}

func NewTransactionHandlerDB(service service.TransactionService) *TransactionHandlerDB {
	return &TransactionHandlerDB{Service: service, Validator: *validator.New()}
}

func (t *TransactionHandlerDB) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, "error", "Method not allowed")
		return
	}

	logger.GetLog().Info().Msg("Creating transaction")
	var req dto.CreateTransactionRequest[domain.Transaction]

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.GetLog().Info().Msg("Invalid request body")
		utils.ErrorResponse(w, http.StatusBadRequest, "error", "Invalid request body")
		return
	}

	/*
	 * Impelemntasi validasi input menggunakan validator di sini
	 * untuk memastikan data input valid dan tidak ada data kosong
	 */
	if err := t.Validator.Struct(req); err != nil {
		errorMessage := utils.CustomValidationError(err)
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, "error", errorMessage)
		return
	}

	trxData := domain.Transaction{
		Type:                 req.Type,
		Nominal:              req.Amount,
		AccountID:            req.AccountID,
		DestinationAccountID: req.DestinationID,
		Note:                 req.Note,
	}

	trx, err := t.Service.CreateTransaction(trxData, req.Amount)
	if err != nil {
		if strings.Contains(err.Error(), "no accounts found") || strings.Contains(err.Error(), "destination account not found") {
			utils.ErrorResponse(w, http.StatusUnprocessableEntity, "error", "Sender id or Receiver id not valid")
		} else if strings.Contains(err.Error(), "insufficient balance") || strings.Contains(err.Error(), "minimum transfer") {
			utils.ErrorResponse(w, http.StatusBadRequest, "error", err.Error())
		} else {
			utils.ErrorResponse(w, http.StatusInternalServerError, "error", "Database error: "+err.Error())
		}
		return
	}

	transaction := dto.TransactionResponse[domain.Transaction]{
		ID:            trx.ID,
		Date:          trx.DateTransaction,
		Type:          trx.Type,
		Amount:        trx.Nominal,
		AccountID:     trx.AccountID,
		Note:          trx.Note,
		DestinationID: trx.DestinationAccountID,
	}

	logger.GetLog().Info().
		Interface("transaction", transaction).
		Msg("Transaction created successfully")
	utils.ResponseJSON(w, transaction, http.StatusCreated, "success", "Transaction created successfully")
}

func (t *TransactionHandlerDB) GetTransactionByAccountID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, "error", "Method not allowed")
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	logger.GetLog().Info().Msg("Get transaction by account ID")
	trx, err := t.Service.GetTransactionByAccountID(id)
	if err != nil {
		if strings.Contains(err.Error(), "no transaction found") {
			utils.ErrorResponse(w, http.StatusNotFound, "error", "Transaction not found for the given account ID")
		} else {
			utils.ErrorResponse(w, http.StatusInternalServerError, "error", err.Error())
		}
		return
	}

	var transactions []dto.TransactionResponse[domain.Transaction]
	for _, transaction := range trx {
		transactions = append(transactions, dto.TransactionResponse[domain.Transaction]{
			ID:            transaction.ID,
			Date:          transaction.DateTransaction,
			Type:          transaction.Type,
			Amount:        transaction.Nominal,
			AccountID:     transaction.AccountID,
			Note:          transaction.Note,
			DestinationID: transaction.DestinationAccountID,
		})
	}

	logger.GetLog().Info().Str("account", id).Str("total", fmt.Sprintf("have: %+v transactions", len(transactions))).Msg("Transaction retrieved successfully")
	utils.ResponseJSON(w, transactions, http.StatusOK, "success", "Transaction retrieved successfully")
}

func (t *TransactionHandlerDB) GetAllTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, "error", "Method not allowed")
		return
	}

	logger.GetLog().Info().Msg("Get all transactions")
	trx, err := t.Service.GetAllTransaction()
	if err != nil && !strings.Contains(err.Error(), "no transaction found") {
		utils.ErrorResponse(w, http.StatusInternalServerError, "error", err.Error())
		return
	}

	var transactions []dto.TransactionResponse[domain.Transaction]
	for _, transaction := range trx {
		transactions = append(transactions, dto.TransactionResponse[domain.Transaction]{
			ID:            transaction.ID,
			Date:          transaction.DateTransaction,
			Type:          transaction.Type,
			Amount:        transaction.Nominal,
			AccountID:     transaction.AccountID,
			Note:          transaction.Note,
			DestinationID: transaction.DestinationAccountID,
		})
	}

	utils.ResponseJSON(w, transactions, http.StatusOK, "success", "Transactions retrieved successfully")
	logger.GetLog().Info().Str("total", fmt.Sprintf("have: %+v transactions", len(transactions))).Msg("Transactions retrieved successfully")
}
