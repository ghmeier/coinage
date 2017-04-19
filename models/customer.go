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
	Quantity  uint64    `json:"quantity" binding:"required"`
}

type Subscribed struct {
	CustomerID  string    `json:"stripeCustomerId"`
	ConnectedID string    `json:"connectedId"`
	RoasterID   uuid.UUID `json:"roasterId"`
	StripeSubID string    `json:"stripeSubId"`
}

/*NewCustomer initializes and returns the id fields of a customer*/
func NewCustomer(userID uuid.UUID, id string) *Customer {
	return &Customer{
		UserID:     userID,
		CustomerID: id,
	}
}

/*NewSubscribeRequest creates a new SubscribeRequest*/
func NewSubscribeRequest(roasterID uuid.UUID, itemID uuid.UUID, frequency Frequency, quantity uint64) *SubscribeRequest {
	return &SubscribeRequest{
		RoasterID: roasterID,
		ItemID:    itemID,
		Frequency: frequency,
		Quantity:  quantity,
	}
}

func NewSubscribed(customerID, connectedID, stripeSubID string, roasterID uuid.UUID) *Subscribed {
	return &Subscribed{
		CustomerID:  customerID,
		ConnectedID: connectedID,
		RoasterID:   roasterID,
		StripeSubID: stripeSubID,
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

func SubscribedFromSQL(rows *sql.Rows) []*Subscribed {
	s := make([]*Subscribed, 0)
	defer rows.Close()

	for rows.Next() {
		sub := &Subscribed{}
		rows.Scan(&sub.CustomerID, &sub.ConnectedID, &sub.RoasterID, &sub.StripeSubID)
		s = append(s, sub)
	}

	return s
}
