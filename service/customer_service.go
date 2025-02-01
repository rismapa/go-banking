package service

import (
	"fmt"

	adapter "github.com/rismapa/go-banking/adapter/repository"
	"github.com/rismapa/go-banking/domain"
)

type CustomerService interface {
	GetAllCustomers() ([]domain.Customer, error)
	AddCustomer(new domain.Customer) ([]domain.Customer, error)
	CreateCustomer(customer domain.Customer) (*domain.Customer, error)
	GetCustomerByID(id string) (*domain.Customer, error)
	UpdateCustomer(id string, customer domain.Customer) (*domain.Customer, error)
}

type CustomerAdapter struct {
	repo adapter.CustomerRepository
}

func NewCustomerService(repository adapter.CustomerRepository) *CustomerAdapter {
	return &CustomerAdapter{repo: repository}
}

// implementasi Primary Port
func (s *CustomerAdapter) GetAllCustomers() ([]domain.Customer, error) {
	return s.repo.FindAll()
}

func (s *CustomerAdapter) AddCustomer(new domain.Customer) ([]domain.Customer, error) {
	return s.repo.AddCustomer(new)
}

func (c *CustomerAdapter) CreateCustomer(customer domain.Customer) (*domain.Customer, error) {
	return c.repo.CreateCustomer(customer)
}

func (s *CustomerAdapter) GetCustomerByID(id string) (*domain.Customer, error) {
	return s.repo.GetCustomerByID(id)
}

func (s *CustomerAdapter) UpdateCustomer(id string, customer domain.Customer) (*domain.Customer, error) {
	_, err := s.repo.GetCustomerByID(id)
	if err != nil {
		return nil, err
	}

	customer.ID = id
	updatedData, err := s.repo.UpdateCustomer(customer)
	if err != nil {
		return nil, fmt.Errorf("internal error while updating customer: %v", err)
	}

	return updatedData, nil
}
