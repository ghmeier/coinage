package gateways

import (
	"fmt"
	"net/http"

	"github.com/pborman/uuid"

	"github.com/ghmeier/bloodlines/config"
	"github.com/ghmeier/bloodlines/gateways"
	"github.com/ghmeier/coinage/models"
)

/*Coinage wraps all API calls to the coinage service*/
type Coinage interface {
	NewRoaster(*models.RoasterRequest) (*models.Roaster, error)
	Roaster(uuid.UUID) (*models.Roaster, error)
	DeleteRoaster(uuid.UUID) error
	Plans(uuid.UUID) ([]*models.Plan, error)
	NewPlan(uuid.UUID, *models.PlanRequest) (*models.Plan, error)
	Plan(uuid.UUID, uuid.UUID) (*models.Plan, error)
	DeletePlan(uuid.UUID, uuid.UUID) error
	NewCustomer(*models.CustomerRequest) (*models.Customer, error)
	Customers(int, int) ([]*models.Customer, error)
	Customer(uuid.UUID) (*models.Customer, error)
	NewSource(uuid.UUID, *models.CustomerRequest) error
	NewSubscription(uuid.UUID, *models.SubscribeRequest) error
	DeleteSubscription(uuid.UUID, string) error
	DeleteCustomer(uuid.UUID) error
}

type coinage struct {
	*gateways.BaseService
	url    string
	client *http.Client
}

/*NewCoinage initializes and returns a Coinage gateway pointing at the host and
  port provided*/
func NewCoinage(config config.Coinage) Coinage {

	var url string
	if config.Port != "" {
		url = fmt.Sprintf("http://%s:%s/api/", config.Host, config.Port)
	} else {
		url = fmt.Sprintf("https://%s/api/", config.Host)
	}

	return &coinage{
		BaseService: gateways.NewBaseService(),
		url:         url,
	}
}

func (c *coinage) NewRoaster(req *models.RoasterRequest) (*models.Roaster, error) {
	url := fmt.Sprintf("%sroaster", c.url)

	var roaster models.Roaster
	err := c.ServiceSend(http.MethodPost, url, req, &roaster)
	if err != nil {
		return nil, err
	}

	return &roaster, nil
}

func (c *coinage) Roaster(id uuid.UUID) (*models.Roaster, error) {
	url := fmt.Sprintf("%sroaster/%s", c.url, id.String())

	var r models.Roaster
	err := c.ServiceSend(http.MethodGet, url, nil, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (c *coinage) DeleteRoaster(id uuid.UUID) error {
	url := fmt.Sprintf("%sroaster/%s", c.url, id.String())

	err := c.ServiceSend(http.MethodDelete, url, nil, nil)
	return err
}

func (c *coinage) Plans(id uuid.UUID) ([]*models.Plan, error) {
	url := fmt.Sprintf("%sroaster/%s/plan", c.url, id.String())

	var r []*models.Plan
	err := c.ServiceSend(http.MethodGet, url, nil, &r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (c *coinage) NewPlan(id uuid.UUID, req *models.PlanRequest) (*models.Plan, error) {
	url := fmt.Sprintf("%sroaster/%s/plan", c.url, id.String())

	var p models.Plan
	err := c.ServiceSend(http.MethodPost, url, req, &p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (c *coinage) Plan(id uuid.UUID, itemID uuid.UUID) (*models.Plan, error) {
	url := fmt.Sprintf("%sroaster/%s/plan/%s", c.url, id.String(), itemID.String())

	var p models.Plan
	err := c.ServiceSend(http.MethodGet, url, nil, &p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (c *coinage) DeletePlan(id uuid.UUID, itemID uuid.UUID) error {
	url := fmt.Sprintf("%sroaster/%s/plan/%s", c.url, id.String(), itemID.String())

	err := c.ServiceSend(http.MethodDelete, url, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *coinage) NewCustomer(req *models.CustomerRequest) (*models.Customer, error) {
	url := fmt.Sprintf("%scustomer", c.url)

	var customer models.Customer
	err := c.ServiceSend(http.MethodPost, url, req, &customer)
	if err != nil {
		return nil, err
	}

	return &customer, nil
}

func (c *coinage) Customers(offset, limit int) ([]*models.Customer, error) {
	url := fmt.Sprintf("%scustomer?offset=%d&limit=%d", c.url, offset, limit)

	var customers []*models.Customer
	err := c.ServiceSend(http.MethodGet, url, nil, customers)
	if err != nil {
		return nil, err
	}

	return customers, nil
}

func (c *coinage) Customer(id uuid.UUID) (*models.Customer, error) {
	url := fmt.Sprintf("%scustomer/%s", c.url, id.String())

	var customer models.Customer
	err := c.ServiceSend(http.MethodGet, url, nil, &customer)
	if err != nil {
		return nil, err
	}

	return &customer, nil
}

func (c *coinage) NewSource(id uuid.UUID, req *models.CustomerRequest) error {
	url := fmt.Sprintf("%scustomer/%s/source", c.url, id.String())

	err := c.ServiceSend(http.MethodPost, url, req, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *coinage) NewSubscription(id uuid.UUID, req *models.SubscribeRequest) error {
	url := fmt.Sprintf("%scustomer/%s/subscription", c.url, id.String())

	err := c.ServiceSend(http.MethodPost, url, req, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *coinage) DeleteSubscription(id uuid.UUID, pid string) error {
	url := fmt.Sprintf("%scustomer/%s/subscription/%s", c.url, id.String(), pid)

	err := c.ServiceSend(http.MethodGet, url, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *coinage) DeleteCustomer(id uuid.UUID) error {
	url := fmt.Sprintf("%scustomer/%s", c.url, id.String())

	err := c.ServiceSend(http.MethodDelete, url, nil, nil)
	if err != nil {
		return err
	}

	return nil
}
