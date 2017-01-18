package models

import (
	"github.com/pborman/uuid"
)

type CustomerAccount struct {
	ID uuid.UUID
}

func NewCustomerAccount() *CustomerAccount {
	return &CustomerAccount{
		ID: uuid.NewUUID(),
	}
}
