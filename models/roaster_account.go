package models

import (
	"database/sql"

	"github.com/pborman/uuid"
)

type RoasterAccount struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"userId"`
	AccountID uuid.UUID `json:"stripeAccountId"`
}

func NewRoasterAccount(userID uuid.UUID, accountID uuid.UUID) *RoasterAccount {
	return &RoasterAccount{
		ID:        uuid.NewUUID(),
		UserID:    userID,
		AccountID: accountID,
	}
}

func FromSql(rows *sql.Rows) ([]*RoasterAccount, error) {
	roasterAccount := make([]*RoasterAccount, 0)

	for rows.Next() {
		c := &RoasterAccount{}
		rows.Scan(&c.ID, c.UserID, c.AccountID)
		roasterAccount = append(roasterAccount, c)
	}

	return roasterAccount, nil
}
