package adapter

import (
	"fmt"

	"github.com/rismapa/go-banking/domain"
)

type CustomerRepositoryMock struct {
	customers []domain.Customer
}

func NewCustomerRepositoryMock() *CustomerRepositoryMock {
	customers := []domain.Customer{
		{
			ID:          "1",
			Name:        "Ari Wibowo",
			City:        "Bandung",
			Zipcode:     "10001",
			DateOfBirth: "1990-01-01",
			Status:      "active",
		},
		{
			ID:          "2",
			Name:        "Budi Santoso",
			City:        "Subang",
			Zipcode:     "90001",
			DateOfBirth: "1985-02-01",
			Status:      "inactive",
		},
	}

	return &CustomerRepositoryMock{
		customers: customers,
	}
}

func (c *CustomerRepositoryMock) FindAll() ([]domain.Customer, error) {
	return c.customers, nil
}

func (c *CustomerRepositoryMock) AddCustomer(new domain.Customer) ([]domain.Customer, error) {
	c.customers = append(c.customers, new)
	return c.customers, nil
}

func (c *CustomerRepositoryMock) CreateCustomer(customer domain.Customer) (*domain.Customer, error) {
	return nil, nil
}

func (c *CustomerRepositoryMock) GetCustomerByID(id string) (*domain.Customer, error) {
	for _, customer := range c.customers {
		if customer.ID == id {
			return &customer, nil
		}
	}

	return nil, nil
}

func (c *CustomerRepositoryMock) UpdateCustomer(updatedCustomer domain.Customer) (*domain.Customer, error) {
	for i, customer := range c.customers {
		if customer.ID == updatedCustomer.ID {
			// Update the fields
			c.customers[i].Name = updatedCustomer.Name
			c.customers[i].City = updatedCustomer.City
			c.customers[i].Zipcode = updatedCustomer.Zipcode
			c.customers[i].DateOfBirth = updatedCustomer.DateOfBirth
			c.customers[i].Status = updatedCustomer.Status

			return &c.customers[i], nil
		}
	}
	return nil, fmt.Errorf("customer with ID %s not found", updatedCustomer.ID)
}
