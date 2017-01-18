package helpers

import (
	"github.com/ghmeier/bloodlines/gateways"
)

type BillingSubscription struct {
	*baseHelper
}

func NewBillingSubscription(sql gateways.SQL) *BillingSubscription {
	return &BillingSubscription{
		baseHelper: &baseHelper{sql: sql},
	}
}
