package helpers

import (
	"github.com/pborman/uuid"

	"github.com/ghmeier/bloodlines/gateways"
	"github.com/jonnykry/coinage/models"
)

type baseHelper struct {
	sql gateways.SQL
}

type Roaster struct {
	*baseHelper
	Stripe gateways.Stripe
}

func NewRoaster(sql gateways.SQL, stripe gateways.Stripe) *Roaster {
	return &Roaster{
		baseHelper: &baseHelper{sql: sql},
		Stripe:     stripe,
	}
}

func (r *Roaster) Insert(req *models.RoasterRequest) (*models.Roaster, error) {
	// TODO: use userID for additional info
	stripeAccount, err := r.Stripe.NewAccount(req.Country)
	if err != nil {
		return "", err
	}

	roaster := models.NewRoaster(req.UserID, stripeAccount.ID)
	err = r.sql.Modify(
		"INSERT INTO roaster_account (id, userId, stripeAccountId)VALUE(?, ?, ?)",
		roaster.ID,
		roaster.UserID,
		roaster.AccountID,
	)
	return roaster, err
}

func (r *Roaster) GetAll(offset int, limit int) ([]*models.Roaster, error) {
	rows, err := r.sql.Select("SELECT id, userId, stripeAccountId FROM roaster_account ORDER BY id ASC LIMIT ?,?", offset, limit)
	if err != nil {
		return nil, err
	}

	roasters, err := models.RoasterFromSql(rows)
	if err != nil {
		return nil, err
	}

	return roasters, nil
}

func (r *Roaster) GetByID(id uuid.UUID) (*models.Roaster, error) {
	rows, err := r.sql.Select("SELECT id, userId, stripeAccountId FROM roaster_account WHERE id=?", id)
	if err != nil {
		return nil, err
	}

	roasters, err := models.RoasterFromSql(rows)
	if err != nil {
		return nil, err
	}

	return roasters[0], nil
}
