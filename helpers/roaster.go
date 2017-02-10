package helpers

import (
	"database/sql"
	"fmt"

	"github.com/pborman/uuid"

	g "github.com/ghmeier/bloodlines/gateways"
	towncenter "github.com/jakelong95/TownCenter/gateways"
	"github.com/jonnykry/coinage/gateways"
	"github.com/jonnykry/coinage/models"
)

type baseHelper struct {
	sql g.SQL
}

type Roaster struct {
	*baseHelper
	Stripe gateways.Stripe
	TC     towncenter.TownCenterI
}

func NewRoaster(sql g.SQL, stripe gateways.Stripe, towncenter towncenter.TownCenterI) *Roaster {
	return &Roaster{
		baseHelper: &baseHelper{sql: sql},
		Stripe:     stripe,
		TC:         towncenter,
	}
}

func (r *Roaster) Insert(req *models.RoasterRequest) (*models.Roaster, error) {
	user, err := r.TC.GetUser(req.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("ERROR: no user for id %s", req.UserID.String())
	} else if user.IsRoaster == 0 {
		return nil, fmt.Errorf("ERROR: no roaster for user %s", req.UserID.String())
	}

	tRoaster, err := r.TC.GetRoaster(user.RoasterId)
	if err != nil {
		return nil, err
	}

	if r == nil {
		return nil, fmt.Errorf("ERROR: no roaster info for id %s", user.RoasterId)
	}

	stripe, err := r.Stripe.NewAccount(req.Country, user, tRoaster)
	if err != nil {
		return nil, err
	}

	roaster := models.NewRoaster(req.UserID, stripe.ID)
	err = r.sql.Modify(
		"INSERT INTO roaster_account (id, userId, stripeAccountId)VALUE(?, ?, ?)",
		roaster.ID,
		roaster.UserID,
		roaster.AccountID,
	)
	if err != nil {
		return nil, err
	}

	roaster.Account = stripe
	return roaster, nil
}

func (r *Roaster) GetByUserID(id uuid.UUID) (*models.Roaster, error) {
	rows, err := r.sql.Select("SELECT id, userId, stripeAccountId FROM roaster_account WHERE userId=?", id)
	if err != nil {
		return nil, err
	}

	return r.account(rows)
}

func (r *Roaster) GetByID(id uuid.UUID) (*models.Roaster, error) {
	rows, err := r.sql.Select("SELECT id, userId, stripeAccountId FROM roaster_account WHERE id=?", id)
	if err != nil {
		return nil, err
	}

	return r.account(rows)
}

func (r *Roaster) account(rows *sql.Rows) (*models.Roaster, error) {
	roasters, _ := models.RoasterFromSql(rows)

	roaster := roasters[0]
	stripe, err := r.Stripe.GetAccount(roaster.AccountID)
	if err != nil {
		return nil, err
	}

	roaster.Account = stripe

	return roaster, nil
}
