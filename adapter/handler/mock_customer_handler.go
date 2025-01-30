package adapter

import (
	"encoding/json"
	"net/http"

	logger "github.com/okyws/go-banking-lib/config"
	"github.com/okyws/go-banking/domain"
	"github.com/okyws/go-banking/service"
)

type CustomerHandler struct {
	service   service.CustomerService
	customers *[]domain.Customer
}

func NewCustomerHandler(service service.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service, customers: &[]domain.Customer{}}
}

func (h *CustomerHandler) GetCustomers(w http.ResponseWriter, r *http.Request) {
	logger.GetLog().Info().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Msg("Get all mock customers")
	customers, err := h.service.GetAllCustomers()
	if err != nil {
		logger.GetLog().Error().Err(err).Msg("Failed to retrieve customers")
		http.Error(w, "Unable to fetch customers", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

func (h *CustomerHandler) AddCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var customer domain.Customer
	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		logger.GetLog().Err(err).Msg("Failed to decode request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.service.AddCustomer(customer)
	if err != nil {
		logger.GetLog().Err(err).Msg("Failed to add customer")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.GetLog().Info().Msg("Customer added successfully")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customer)
}
