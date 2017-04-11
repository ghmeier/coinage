package gateways

import (
	"fmt"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"

	"github.com/ghmeier/bloodlines/config"
	"github.com/ghmeier/coinage/models"
	tmodels "github.com/jakelong95/TownCenter/models"
	item "github.com/lcollin/warehouse/models"
)

/*Stripe wraps the stripe-go api for coinage use*/
type Stripe interface {
	NewCustomer(token, userID string) (string, error)
	GetCustomer(id string) (*stripe.Customer, error)
	DeleteCustomer(id string) error
	AddSource(id, token string) (*stripe.Customer, error)
	NewAccount(user *tmodels.User, roaster *tmodels.Roaster) (*stripe.Account, error)
	NewPlan(secret string, item *item.Item, freq models.Frequency) (*stripe.Plan, error)
	GetPlan(secret, pid string) (*stripe.Plan, error)
	Subscribe(secret, id, planID string) (*stripe.Sub, error)
}

type stripeS struct {
	config config.Stripe
	c      *client.API
}

/*NewStripe initializes and returns a new Stripe implementation configured
  by the provided config*/
func NewStripe(config config.Stripe) Stripe {
	s := &stripeS{
		config: config,
		c:      client.New(config.Secret, nil),
	}

	return s
}

func (s *stripeS) NewCustomer(token string, userID string) (string, error) {
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
func (s *stripeS) GetCustomer(id string) (*stripe.Customer, error) {
	c, err := s.c.Customers.Get(id, nil)
	if err != nil {
		return nil, err
	}

	return c, nil
}

/*DeleteCustomer removes customer from stripe by customer id*/
func (s *stripeS) DeleteCustomer(id string) error {
	_, err := s.c.Customers.Del(id)
	return err
}

/*AddSource creates and adds a new Default Source for the customer*/
func (s *stripeS) AddSource(id string, token string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Source: &stripe.SourceParams{Token: token},
	}

	c, err := s.c.Customers.Update(id, params)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (s *stripeS) NewAccount(user *tmodels.User, roaster *tmodels.Roaster) (*stripe.Account, error) {
	params := &stripe.AccountParams{
		Managed:      true,
		Country:      "US",
		Email:        roaster.Email,
		BusinessName: roaster.Name,
		LegalEntity: &stripe.LegalEntity{
			BusinessName: roaster.Name,
			Address: stripe.Address{
				City:    roaster.AddressCity,
				Country: "US",
				Line1:   roaster.AddressLine1,
				Line2:   roaster.AddressLine2,
				Zip:     roaster.AddressZip,
				State:   roaster.AddressState,
			},
			PersonalAddress: stripe.Address{
				City:    user.AddressCity,
				Country: "US",
				Line1:   user.AddressLine1,
				Line2:   user.AddressLine2,
				Zip:     user.AddressZip,
				State:   user.AddressState,
			},
			First:       user.FirstName,
			Last:        user.LastName,
			PhoneNumber: roaster.Phone,
			Type:        "company",
		},
	}

	account, err := s.c.Account.New(params)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *stripeS) NewPlan(secret string, item *item.Item, freq models.Frequency) (*stripe.Plan, error) {
	interval, ok := models.ToFrequency(freq)
	if !ok {
		return nil, fmt.Errorf("ERROR: no frequency for interval %s", freq)
	}

	client := client.New(secret, nil)
	params := &stripe.PlanParams{
		ID:            fmt.Sprintf("%s-%s", item.ID.String(), string(freq)),
		Amount:        uint64(item.ConsumerPrice * 100),
		Currency:      "usd",
		Interval:      "week",
		IntervalCount: uint64(interval),
		Name:          fmt.Sprintf("%s %s", item.Name, string(freq)),
		Statement:     fmt.Sprintf("exp-%.18s", item.Name),
	}
	params.AddMeta("itemId", item.ID.String())
	plan, err := client.Plans.New(params)

	if err != nil {
		return nil, err
	}

	return plan, nil
}

func (s *stripeS) GetPlan(secret string, pid string) (*stripe.Plan, error) {
	client := client.New(secret, nil)
	plan, err := client.Plans.Get(pid, nil)
	if err != nil {
		return nil, err
	}

	return plan, nil
}

func (s *stripeS) Subscribe(secret, customerID, planID string) (*stripe.Sub, error) {
	client := client.New(secret, nil)
	sub, err := client.Subs.New(&stripe.SubParams{Customer: customerID, Plan: planID})
	if err != nil {
		return nil, err
	}

	return sub, nil
}
