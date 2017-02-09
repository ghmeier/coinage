package helpers

import (
	"fmt"

	"github.com/pborman/uuid"

	g "github.com/ghmeier/bloodlines/gateways"
	towncenter "github.com/jakelong95/TownCenter/gateways"
	"github.com/jonnykry/coinage/gateways"
	"github.com/jonnykry/coinage/models"
)

type Customer struct {
	*baseHelper
	Stripe gateways.Stripe
	TC     towncenter.TownCenterI
}

func NewCustomer(sql g.SQL, stripe gateways.Stripe, tc towncenter.TownCenterI) *Customer {
	return &Customer{
		baseHelper: &baseHelper{sql: sql},
		Stripe:     stripe,
		TC:         tc,
	}
}

func (c *Customer) Insert(req *models.CustomerRequest) (*models.Customer, error) {
	user, err := c.TC.GetUser(req.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("ERROR: no user for id %s", req.UserID.String())
	}

	customerID, err := c.Stripe.NewCustomer(req.Token, req.UserID.String())
	if err != nil {
		return nil, err
	}

	customer := models.NewCustomer(req.UserID, customerID)
	err = c.sql.Modify("INSERT INTO customer_account (id, userId, stripeCustomerId)VALUES(?,?,?)",
		customer.ID,
		customer.UserID,
		customer.CustomerID)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (c *Customer) View(id uuid.UUID) (*models.Customer, error) {
	rows, err := c.sql.Select("SELECT id, userId, stripeCustomerId FROM customer_account WHERE userId=?", id.String())
	if err != nil {
		return nil, err
	}

	// no error possible right now
	customers, _ := models.CustomersFromSQL(rows)
	customer := customers[0]

	stripeCustomer, err := c.Stripe.GetCustomer(customer.CustomerID)
	if err != nil {
		return nil, err
	}

	customer.Sources = stripeCustomer.Sources
	customer.Subscriptions = stripeCustomer.Subs
	customer.Meta = stripeCustomer.Meta
	return customer, nil
}

func (c *Customer) AddSource(id uuid.UUID, token string) error {
	customer, err := c.View(id)
	if err != nil {
		return err
	}

	_, err = c.Stripe.AddSource(customer.CustomerID, token)
	return err
}

func (c *Customer) Delete(id uuid.UUID) error {
	err := c.sql.Modify("DELETE FROM customer_account WHERE userId=?", id)
	return err
}
