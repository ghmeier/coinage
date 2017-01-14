package helpers

import (
	"github.com/ghmeier/bloodlines/gateways"
)

type CustomerAccount struct {
	*baseHelper
}

func NewCustomerAccount(sql gateways.SQL) *CustomerAccount {
	return &CustomerAccount{baseHelper: &baseHelper{sql: sql}}
}
