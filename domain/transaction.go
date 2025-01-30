package domain

type Transaction struct {
	ID                   string  `db:"id"`
	DateTransaction      string  `db:"date_transaction" validate:"required,datetime=2006-01-02 15:04:05"`
	Type                 string  `db:"type" validate:"required,oneof=debit credit transfer"`
	Nominal              float64 `db:"nominal" validate:"required,gt=10000"`
	AccountID            string  `db:"account_id" validate:"required,uuid"`
	Note                 string  `db:"note"  validate:"omitempty,max=255"`
	DestinationAccountID string  `db:"destination_account_id" validate:"omitempty,uuid"`
}
