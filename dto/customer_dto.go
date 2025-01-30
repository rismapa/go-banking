package dto

type CustomerRequestDTO[T any] struct {
	ID          string `json:"cust_id"`
	Name        string `json:"cust_name" validate:"required,min=3,max=100"`
	City        string `json:"cust_city" validate:"required,max=100"`
	Zipcode     string `json:"cust_zipcode" validate:"required,max=10"`
	DateOfBirth string `json:"cust_dob" validate:"required,datetime=2006-01-02"`
	Status      string `json:"cust_status" validate:"required,oneof=active inactive"`
}

type CustomerResponseDTO[T any] struct {
	ID      string `json:"cust_id"`
	Name    string `json:"cust_name"`
	City    string `json:"cust_city"`
	Zipcode string `json:"cust_zipcode"`
	Status  string `json:"cust_status"`
}
