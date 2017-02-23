package models

import (
	"database/sql"

	"github.com/pborman/uuid"
	"github.com/stripe/stripe-go"
)

/*Roaster has information retrieved from stripe and the db
  about billing for roaster entities*/
type Roaster struct {
	//ID is the roaster ID in towncenter
	ID        uuid.UUID       `json:"id"`
	AccountID string          `json:"stripeAccountId"`
	Account   *stripe.Account `json:"account"`
}

/*RoasterRequest has information used in creating a roaster
  managed account in stripe*/
type RoasterRequest struct {
	UserID  uuid.UUID `json:"userId" binding:"required"`
	Country string    `json:"country" binding:"required"`
	/* TODO: more info as we need it */
}

/*NewRoaster initialized and returns a roaster model*/
func NewRoaster(id uuid.UUID, accountID string) *Roaster {
	return &Roaster{
		ID:        uuid.NewUUID(),
		AccountID: accountID,
	}
}

/*RoasterFromSQL maps an sql row to roaster properties,
  where order matters*/
func RoasterFromSQL(rows *sql.Rows) ([]*Roaster, error) {
	roasters := make([]*Roaster, 0)

	for rows.Next() {
		c := &Roaster{}
		rows.Scan(&c.ID, &c.AccountID)
		roasters = append(roasters, c)
	}

	return roasters, nil
}
