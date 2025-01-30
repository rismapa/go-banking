package dto

type CreateTransactionRequest[T any] struct {
	Type          string  `json:"trx_type" validate:"required,oneof=debit credit transfer"`
	Amount        float64 `json:"trx_amount" validate:"required,gt=0"`
	AccountID     string  `json:"trx_account_id" validate:"required,uuid"`
	Note          string  `json:"trx_note,omitempty" validate:"omitempty,max=1000"`
	DestinationID string  `json:"trx_destination_id,omitempty" validate:"omitempty,uuid,required_if=Type transfer"`
}

type TransactionResponse[T any] struct {
	ID            string  `json:"trx_id"`
	Date          string  `json:"trx_date"`
	Type          string  `json:"trx_type"`
	Amount        float64 `json:"trx_amount"`
	AccountID     string  `json:"trx_account_id"`
	Note          string  `json:"trx_note"`
	DestinationID string  `json:"trx_destination_id,omitempty"`
}
