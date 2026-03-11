package storage

import (
	"Customer/internal/models/customer"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrUrlNotFound = errors.New("url not found")
	ErrUrlExist    = errors.New("url exist")
)

type OrderService interface {
	AddURL(customer customer.Customer) (uuid.UUID, error)

	DeleteURL(id uuid.UUID) error

	GetAllURL() ([]customer.Customer, error)

	GetByIdURL(id uuid.UUID) (uuid.UUID, error)

	UpdateURL(customer customer.Customer) error

	IsCustomerCreatedURL(id uuid.UUID) (bool, error)
}
