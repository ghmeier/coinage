package models

import (
	"database/sql"
	"github.com/stripe/stripe-go"

	"github.com/pborman/uuid"
)

type Customer struct {
	ID            uuid.UUID          `json:"id"`
	UserID        uuid.UUID          `json:"userId"`
	CustomerID    string             `json:"stripeCustomerId"`
	Subscriptions *stripe.SubList    `json:"subscriptions"`
	Sources       *stripe.SourceList `json:"sources"`
	Meta          map[string]string  `json:"metadata"`
	// SourceID       string    `json:"stripeCardId"`
	// SubscriptionID string `json:"stripeSubscriptionId"`
	// PlanID         string `json:"stripePlanId"`
}

type CustomerRequest struct {
	UserID uuid.UUID `json:"userId" binding:"required"`
	Token  string    `json:"token" binding:"required"`
}

type SubscribeRequest struct {
	RoasterID uuid.UUID `json:"roasterId" binding:"required"`
	ItemID    uuid.UUID `json:"itemId" binding:"required"`
	Frequency Frequency `json:"frequency" binding:"required"`
}

func NewCustomer(userID uuid.UUID, id string) *Customer {
	return &Customer{
		ID:         uuid.NewUUID(),
		UserID:     userID,
		CustomerID: id,
	}
}

func CustomersFromSQL(rows *sql.Rows) ([]*Customer, error) {
	customers := make([]*Customer, 0)
	defer rows.Close()

	for rows.Next() {
		c := &Customer{}

		rows.Scan(&c.ID, &c.UserID, &c.CustomerID)

		customers = append(customers, c)
	}

	return customers, nil
}
