package helpers

import (
	"database/sql"
	"fmt"

	"github.com/pborman/uuid"

	g "github.com/ghmeier/bloodlines/gateways"
	"github.com/ghmeier/coinage/gateways"
	"github.com/ghmeier/coinage/models"
	towncenter "github.com/jakelong95/TownCenter/gateways"
	t "github.com/jakelong95/TownCenter/models"
)

type base struct {
	sql g.SQL
}

/*Roaster helps with manipulating roaster properties*/
type Roaster struct {
	*base
	Stripe gateways.Stripe
	TC     towncenter.TownCenterI
}

/*NewRoaster initializes and returns a roaster with the given gateways*/
func NewRoaster(sql g.SQL, stripe gateways.Stripe, towncenter towncenter.TownCenterI) *Roaster {
	return &Roaster{
		base:   &base{sql: sql},
		Stripe: stripe,
		TC:     towncenter,
	}
}

/*Insert creates a new roaster's stripe information and adds a record to the db*/
func (r *Roaster) Insert(req *models.RoasterRequest) (*models.Roaster, error) {
	user, tRoaster, err := r.roaster(req.UserID)
	if err != nil {
		return nil, err
	}
	stripe, err := r.Stripe.NewAccount(req.Country, user, tRoaster)
	if err != nil {
		return nil, err
	}

	roaster := models.NewRoaster(tRoaster.ID, stripe.ID)
	err = r.sql.Modify(
		"INSERT INTO roaster_account (id, stripeAccountId)VALUE(?, ?, ?)",
		roaster.ID,
		roaster.AccountID,
	)
	if err != nil {
		return nil, err
	}

	roaster.Account = stripe
	return roaster, nil
}

/*GetByUserID returns the roaster account associated with a user id*/
func (r *Roaster) GetByUserID(id uuid.UUID) (*models.Roaster, error) {
	_, roaster, err := r.roaster(id)
	if err != nil {
		return nil, err
	}

	return r.Get(roaster.ID)
}

/*Get returns the roaster account associated with the given id*/
func (r *Roaster) Get(id uuid.UUID) (*models.Roaster, error) {
	rows, err := r.sql.Select("SELECT id, stripeAccountId FROM roaster_account WHERE id=?", id)
	if err != nil {
		return nil, err
	}

	return r.account(rows)
}

func (r *Roaster) account(rows *sql.Rows) (*models.Roaster, error) {
	roasters, _ := models.RoasterFromSQL(rows)
	if len(roasters) < 1 {
		return nil, nil
	}

	roaster := roasters[0]
	stripe, err := r.Stripe.GetAccount(roaster.AccountID)
	if err != nil {
		return nil, err
	}

	roaster.Account = stripe

	return roaster, nil
}

/* roaster returns a towncenter user && roaster by user id. errors otherwise */
func (r *Roaster) roaster(id uuid.UUID) (*t.User, *t.Roaster, error) {
	u, err := r.TC.GetUser(id)
	if err != nil {
		return nil, nil, err
	}

	if u == nil {
		return nil, nil, fmt.Errorf("ERROR: no user for id %s", id.String())
	} else if u.RoasterId == nil {
		return nil, nil, fmt.Errorf("ERROR: no roaster for user %s", id.String())
	}

	roaster, err := r.TC.GetRoaster(u.RoasterId)
	if err != nil {
		return nil, nil, err
	}

	if roaster == nil {
		return nil, nil, fmt.Errorf("ERROR: no roaster info for id %s", u.RoasterId.String())
	}

	return u, roaster, nil
}
