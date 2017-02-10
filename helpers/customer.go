package helpers

import (
	"fmt"

	"github.com/pborman/uuid"

	g "github.com/ghmeier/bloodlines/gateways"
	t "github.com/jakelong95/TownCenter/gateways"
	"github.com/jonnykry/coinage/gateways"
	"github.com/jonnykry/coinage/models"
	c "github.com/yuderekyu/covenant/gateways"
	//sub "github.com/yuderekyu/covenant/models"
)

type Customer struct {
	*baseHelper
	Stripe   gateways.Stripe
	Covenant c.Covenant
	TC       t.TownCenterI
}

func NewCustomer(sql g.SQL, stripe gateways.Stripe, tc t.TownCenterI, cov c.Covenant) *Customer {
	return &Customer{
		baseHelper: &baseHelper{sql: sql},
		Covenant:   cov,
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

// userID uuid.UUID, createdAt string, startAt string, shopID uuid.UUID, ozInBag float64,
// beanName string, roastName string, price float64
func (c *Customer) Subscribe(id uuid.UUID, plan *models.Plan, freq models.Frequency) error {
	customer, err := c.View(id)
	if err != nil {
		return err
	}

	interval, ok := models.ToFrequency(string(freq))
	if !ok {
		return fmt.Errorf("ERROR: invalid frequency %s", freq)
	}

	stripe := plan.PlanIDs[interval]

	_, err = c.Stripe.Subscribe(customer.CustomerID, stripe)
	if err != nil {
		return err
	}

	// TODO: add subs to covenant
	// subscription := sub.NewSubscription(id, time.Now(), time.Now(), plan.RoasterID)
	// c.Covenant.NewSubscription()

	return err
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
	customer, err := c.View(id)
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
