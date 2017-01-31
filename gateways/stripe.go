package gateways

import (
	"fmt"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"

	"github.com/ghmeier/bloodlines/config"
)

type Stripe interface {
	NewCustomer(token string, userID string) (string, error)
	GetCustomer(id string) (*stripe.Customer, error)
	DeleteCustomer(id string) error
	AddSource(id string, token string) (*stripe.Customer, error)
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
	stripe.Key = config.Secret

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
