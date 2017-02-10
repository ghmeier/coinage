package models

import (
	"database/sql"

	"github.com/pborman/uuid"
	"github.com/stripe/stripe-go"
)

type Roaster struct {
	//ID is the roaster ID in towncenter
	ID        uuid.UUID       `json:"id"`
	AccountID string          `json:"stripeAccountId"`
	Account   *stripe.Account `json:"account"`
}

type RoasterRequest struct {
	UserID  uuid.UUID `json:"userId" binding:"required"`
	Country string    `json:"country" binding:"required"`
	/* TODO: more info as we need it */
}

func NewRoaster(id uuid.UUID, accountID string) *Roaster {
	return &Roaster{
		ID:        uuid.NewUUID(),
		AccountID: accountID,
	}
}

func RoasterFromSql(rows *sql.Rows) ([]*Roaster, error) {
	roasters := make([]*Roaster, 0)

	for rows.Next() {
		c := &Roaster{}
		rows.Scan(&c.ID, &c.AccountID)
		roasters = append(roasters, c)
	}

	return roasters, nil
}
