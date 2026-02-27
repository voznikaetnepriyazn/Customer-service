package customer

import (
	"github.com/google/uuid"
)

type Customer struct {
	Id          uuid.UUID
	Name        string
	Email       string
	City        string
	FullAdress  string
	PostalCode  int64
	Phonenumber int64
}
