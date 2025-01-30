package domain

type Customer struct {
	ID          string `json:"cust_id" db:"id"`
	Name        string `json:"cust_name" db:"name" validate:"required,min=3,max=100"`
	City        string `json:"cust_city" db:"city" validate:"required,max=100"`
	Zipcode     string `json:"cust_zipcode" db:"zipcode" validate:"required,len=10"`
	DateOfBirth string `json:"cust_dob" db:"date_of_birth" validate:"required,datetime=2006-01-02"`
	Status      string `json:"cust_status" db:"status" validate:"required,oneof=active inactive"`
}
