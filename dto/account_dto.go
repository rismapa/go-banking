package dto

import "github.com/rismapa/go-banking/domain"

type AccountRequest[T any] struct {
	Customer_ID string  `json:"cust_id" validate:"required,uuid"`
	Username    string  `json:"acc_username" validate:"required,min=3,max=100"`
	Password    string  `json:"acc_password" validate:"required,min=8,max=100"`
	Balance     float64 `json:"acc_balance" validate:"gte=0,min=0"`
	Currency    string  `json:"acc_currency" validate:"required,len=3"`
	Status      bool    `json:"acc_status" validate:"boolean"`
}

type AccountResponse[T any] struct {
	ID          string  `json:"acc_id" db:"id"`
	Customer_ID string  `json:"cust_id" db:"customer_id" validate:"required,uuid"`
	Username    string  `json:"acc_username" db:"username" validate:"required,min=3,max=100"`
	Balance     float64 `json:"acc_balance" db:"balance" validate:"gte=0,min=0"`
	Currency    string  `json:"acc_currency" db:"currency" validate:"required,len=3"`
	Status      bool    `json:"acc_status" db:"status" validate:"boolean"`
}

type AccountWithCustomer struct {
	ID            string  `json:"acc_id" db:"id"`
	Customer_Name string  `json:"cust_name" db:"name"`
	Username      string  `json:"acc_username" db:"username"`
	Balance       float64 `json:"acc_balance"`
	Currency      string  `json:"acc_currency"`
	Status        bool    `json:"acc_status"`
}

type AccountByCustomerIDResponse struct {
	CustomerData *domain.Customer `json:"cus_data"`
	AccountData  []domain.Account `json:"acc_data"`
}
