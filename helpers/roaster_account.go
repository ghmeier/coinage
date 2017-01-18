package helpers

import (
	"github.com/pborman/uuid"

	"github.com/ghmeier/bloodlines/gateways"
	"github.com/jonnykry/coinage/models"
)

type baseHelper struct {
	sql gateways.SQL
	//stats *statsd.Client
}

type RoasterAccount struct {
	*baseHelper
}

func NewRoasterAccount(sql gateways.SQL) *RoasterAccount {
	return &RoasterAccount{baseHelper: &baseHelper{sql: sql}}
}

func (r *RoasterAccount) Insert(account *models.RoasterAccount) error {
	err := r.sql.Modify(
		"INSERT INTO roaster_account (id, userId, stripeAccountId)VALUE(?, ?, ?)",
		account.ID,
		account.UserID,
		account.AccountID,
	)
	return err
}

func (r *RoasterAccount) GetAll(offset int, limit int) ([]*models.RoasterAccount, error) {
	rows, err := r.sql.Select("SELECT id, userId, stripeAccountId from roaster_account ORDER BY id ASC LIMIT ?,?", offset, limit)
	if err != nil {
		return nil, err
	}

	accounts, err := models.FromSql(rows)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (r *RoasterAccount) GetByID(id uuid.UUID) (*models.RoasterAccount, error) {
	rows, err := r.sql.Select("SELECT id, userId, stripeAccountId FROM roaster_account WHERE id=?", id)
	if err != nil {
		return nil, err
	}

	accounts, err := models.FromSql(rows)
	if err != nil {
		return nil, err
	}

	return accounts[0], nil
}

func (r *RoasterAccount) Update(account *models.RoasterAccount) error {
	err := r.sql.Modify("UPDATE roaster_account SET userId=?,stripeAccountId=? WHERE id=?",
		account.UserID,
		account.AccountID,
		account.ID,
	)
	return err
}
