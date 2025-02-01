package adapter

import (
	"fmt"

	"github.com/rismapa/go-banking/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CustomerRepository interface {
	FindAll() ([]domain.Customer, error)
	AddCustomer(new domain.Customer) ([]domain.Customer, error)
	CreateCustomer(customer domain.Customer) (*domain.Customer, error)
	GetCustomerByID(id string) (*domain.Customer, error)
	UpdateCustomer(customer domain.Customer) (*domain.Customer, error)
}

type CustomerRepositoryDB struct {
	DB *sqlx.DB
}

func NewCustomerRepositoryDB(db *sqlx.DB) *CustomerRepositoryDB {
	return &CustomerRepositoryDB{DB: db}
}

func (r *CustomerRepositoryDB) FindAll() ([]domain.Customer, error) {
	var customers []domain.Customer
	err := r.DB.Select(&customers, "SELECT id, name, city, zipcode, date_of_birth, status FROM customers")
	if err != nil {
		return nil, fmt.Errorf("failed to get customers: %v", err)
	}

	if len(customers) == 0 {
		return nil, fmt.Errorf("no customers found")
	}

	return customers, nil
}

func (r *CustomerRepositoryDB) CreateCustomer(customer domain.Customer) (*domain.Customer, error) {
	customer.ID = uuid.New().String()
	_, err := r.DB.Exec("INSERT INTO customers (id, name, city, zipcode, date_of_birth, status) VALUES (?, ?, ?, ?, ?, ?)",
		customer.ID, customer.Name, customer.City, customer.Zipcode, customer.DateOfBirth, customer.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %v", err)
	}

	return &customer, nil
}

func (r *CustomerRepositoryDB) GetCustomerByID(id string) (*domain.Customer, error) {
	var customer domain.Customer
	query := "SELECT id, name, city, zipcode, date_of_birth, status FROM customers WHERE id = ?"
	err := r.DB.Get(&customer, query, id)
	if err != nil {
		if customer == (domain.Customer{}) {
			return nil, fmt.Errorf("no customers found")
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	return &customer, nil
}

func (r *CustomerRepositoryDB) UpdateCustomer(customer domain.Customer) (*domain.Customer, error) {
	query := "UPDATE customers SET name = ?, city = ?, zipcode = ?, date_of_birth = ?, status = ? WHERE id = ?"
	_, err := r.DB.Exec(query, customer.Name, customer.City, customer.Zipcode, customer.DateOfBirth, customer.Status, customer.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update customer: %v", err)
	}

	return &customer, nil
}

func (c *CustomerRepositoryDB) AddCustomer(new domain.Customer) ([]domain.Customer, error) {
	// tidak ada implementasi karena ini untuk mock data
	return nil, nil
}
