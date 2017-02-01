package gateways

import (
	"fmt"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"

	"github.com/ghmeier/bloodlines/config"
	"github.com/jonnykry/coinage/models"
)

type Stripe interface {
	NewCustomer(token string, userID string) (string, error)
	GetCustomer(id string) (*stripe.Customer, error)
	DeleteCustomer(id string) error
	AddSource(id string, token string) (*stripe.Customer, error)
	NewAccount(country string) (*stripe.Account, error)
	GetAccount(id string) (*stripe.Account, error)
	NewPlan(id string, req *models.PlanRequest) (*stripe.Plan, error)
	GetPlan(id string, pid string) (*stripe.Plan, error)
}

type StripeS struct {
	config config.Stripe
	c      *client.API
}

func NewStripe(config config.Stripe) Stripe {
	s := &StripeS{
		config: config,
		c:      client.New(config.Secret, nil),
	}

	return s
}

func (s *StripeS) NewCustomer(token string, userID string) (string, error) {
	params := &stripe.CustomerParams{
		Desc: fmt.Sprintf("Customer for user: %s", userID),
	}

	params.SetSource(token)
	c, err := s.c.Customers.New(params)
	if err != nil {
		return "", err
	}
	return c.ID, err
}

/*GetCustomer returns a stripe customer by their cutsomerID*/
func (s *StripeS) GetCustomer(id string) (*stripe.Customer, error) {
	c, err := s.c.Customers.Get(id, nil)
	if err != nil {
		return nil, err
	}

	return c, nil
}

/*DeleteCustomer removes customer from stripe by customer id*/
func (s *StripeS) DeleteCustomer(id string) error {
	_, err := s.c.Customers.Del(id)
	return err
}

/*AddSource creates and adds a new Default Source for the customer*/
func (s *StripeS) AddSource(id string, token string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Source: &stripe.SourceParams{Token: token},
	}

	c, err := s.c.Customers.Update(id, params)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (s *StripeS) NewAccount(country string) (*stripe.Account, error) {
	params := &stripe.AccountParams{
		Managed: true,
		Country: country,
	}

	account, err := s.c.Account.New(params)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *StripeS) GetAccount(id string) (*stripe.Account, error) {
	account, err := s.c.Account.GetByID(id, nil)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *StripeS) NewPlan(id string, req *models.PlanRequest) (*stripe.Plan, error) {
	account, err := s.GetAccount(id)
	if err != nil {
		return nil, err
	}

	client := client.New(account.Keys.Secret, nil)
	plan, err := client.Plans.New(&stripe.PlanParams{})
	if err != nil {
		return nil, err
	}

	return plan, nil
}

func (s *StripeS) GetPlan(id string, pid string) (*stripe.Plan, error) {
	account, err := s.GetAccount(id)
	if err != nil {
		return nil, err
	}

	client := client.New(account.Keys.Secret, nil)
	plan, err := client.Plans.Get(id, nil)
	if err != nil {
		return nil, err
	}

	return plan, nil
}
