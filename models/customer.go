package models

import (
	"database/sql"

	"github.com/pborman/uuid"
)

/*Customer is the stripe customer data connected to a userID*/
type Customer struct {
	UserID     uuid.UUID `json:"userId"`
	CustomerID string    `json:"stripeCustomerId"`
	// SourceID       string    `json:"stripeCardId"`
	// SubscriptionID string `json:"stripeSubscriptionId"`
	// PlanID         string `json:"stripePlanId"`
}

/*CustomerRequest is the information needed to create/update a stripe customer*/
type CustomerRequest struct {
	UserID uuid.UUID `json:"userId" binding:"required"`
	Token  string    `json:"token" binding:"required"`
}

/*SubscribeRequest is the information needed to subscribe a customer
  to a roaster plan*/
type SubscribeRequest struct {
	RoasterID uuid.UUID `json:"roasterId" binding:"required"`
	ItemID    uuid.UUID `json:"itemId" binding:"required"`
	Frequency Frequency `json:"frequency" binding:"required"`
}

/*NewCustomer initializes and returns the id fields of a customer*/
func NewCustomer(userID uuid.UUID, id string) *Customer {
	return &Customer{
		UserID:     userID,
		CustomerID: id,
	}
}

/*NewSubscribeRequest creates a new SubscribeRequest*/
func NewSubscribeRequest(roasterID uuid.UUID, itemID uuid.UUID, frequency Frequency) *SubscribeRequest {
	return &SubscribeRequest{
		RoasterID: roasterID,
		ItemID:    itemID,
		Frequency: frequency,
	}
}

/*CustomersFromSQL returns a customer model slice from sql rows*/
func CustomersFromSQL(rows *sql.Rows) ([]*Customer, error) {
	customers := make([]*Customer, 0)
	defer rows.Close()

	for rows.Next() {
		c := &Customer{}

		rows.Scan(&c.UserID, &c.CustomerID)

		customers = append(customers, c)
	}

	return customers, nil
}
