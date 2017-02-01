package models

import (
	"database/sql"

	"github.com/pborman/uuid"
)

type Roaster struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"userId"`
	AccountID uuid.UUID `json:"stripeAccountId"`
}

type RoasterRequest struct {
	UserID  uuid.UUID `json:"userId" binding:"required"`
	Country string    `json:"country" binding:"required"`
	/* TODO: more info as we need it */
}

func NewRoaster(userID uuid.UUID, accountID uuid.UUID) *Roaster {
	return &Roaster{
		ID:        uuid.NewUUID(),
		UserID:    userID,
		AccountID: accountID,
	}
}

func RoasterFromSql(rows *sql.Rows) ([]*Roaster, error) {
	roasters := make([]*Roaster, 0)

	for rows.Next() {
		c := &Roaster{}
		rows.Scan(&c.ID, &c.UserID, &c.AccountID)
		roasters = append(roasters, c)
	}

	return roaster, nil
}
