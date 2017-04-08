package helpers

import (
	"fmt"

	"github.com/pborman/uuid"

	g "github.com/ghmeier/bloodlines/gateways"
	"github.com/ghmeier/coinage/gateways"
	"github.com/ghmeier/coinage/models"
	t "github.com/jakelong95/TownCenter/gateways"
)

/*Customer helps with creating and manipulating stripe customers*/
type Customer struct {
	*base
	Stripe gateways.Stripe
	TC     t.TownCenterI
}

/*NewCustomer initializes a Customer with the given gateways*/
func NewCustomer(sql g.SQL, stripe gateways.Stripe, tc t.TownCenterI) *Customer {
	return &Customer{
		base:   &base{sql: sql},
		Stripe: stripe,
		TC:     tc,
	}
}

/*Insert creates a new stripe customer with the given id and token, inserting a record
  into the db*/
func (c *Customer) Insert(req *models.CustomerRequest) (*models.Customer, error) {
	customer, err := c.Get(req.UserID)
	if err != nil {
		return nil, err
	}
	if customer != nil {
		return c.AddSource(customer, req.Token)
	}

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

	customer = models.NewCustomer(req.UserID, customerID)
	err = c.sql.Modify("INSERT INTO customer_account (userId, stripeCustomerId)VALUES(?,?,?)",
		customer.UserID,
		customer.CustomerID)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

/*Get returns a customer associated with the given UserID*/
func (c *Customer) Get(id uuid.UUID) (*models.Customer, error) {
	rows, err := c.sql.Select("SELECT userId, stripeCustomerId FROM customer_account WHERE userId=?", id.String())
	if err != nil {
		return nil, err
	}

	// no error possible right now
	customers, _ := models.CustomersFromSQL(rows)
	if len(customers) < 1 {
		// no customers for that userID
		return nil, nil
	}

	customer := customers[0]
	if err != nil {
		return nil, err
	}

	return customer, nil
}

/*Subscribe creates a new subscription to the provided plan at the given Frequency in stripe*/
func (c *Customer) Subscribe(id uuid.UUID, plan *models.Plan, freq models.Frequency) error {
	customer, err := c.Get(id)
	if err != nil {
		return err
	}

	interval, ok := models.ToFrequency(freq)
	if !ok {
		return fmt.Errorf("ERROR: invalid frequency %s", freq)
	}

	stripe := plan.PlanIDs[interval-1]

	_, err = c.Stripe.Subscribe(customer.CustomerID, stripe)
	return err
}

/*AddSource creates a new stripe source and sets it as default for the
  given customer*/
func (c *Customer) AddSource(customer *models.Customer, token string) (*models.Customer, error) {
	_, err := c.Stripe.AddSource(customer.CustomerID, token)
	return customer, err
}

/*Delete removes a customer from strip and the db*/
func (c *Customer) Delete(id uuid.UUID) error {
	customer, err := c.Get(id)
	if err != nil {
		return err
	}

	err = c.Stripe.DeleteCustomer(customer.CustomerID)
	if err != nil {
		return err
	}

	err = c.sql.Modify("DELETE FROM customer_account WHERE userId=?", id)
	return err
}
