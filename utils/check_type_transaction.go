package utils

// IsValidTransactionType validates if the transaction type is one of the allowed values.
func IsValidTransactionType(transactionType string) bool {
	validTypes := map[string]bool{
		"transfer": true,
		"credit":   true,
		"debit":    true,
	}

	return validTypes[transactionType]
}
