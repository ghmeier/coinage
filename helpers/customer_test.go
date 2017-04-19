package helpers

import (
	//"database/sql"
	"fmt"
	"testing"

	bgateways "github.com/ghmeier/bloodlines/gateways"
	mocks "github.com/ghmeier/coinage/_mocks/gateways"
	"github.com/ghmeier/coinage/models"
	tmocks "github.com/jakelong95/TownCenter/_mocks"
	tmodels "github.com/jakelong95/TownCenter/models"
	wmocks "github.com/lcollin/warehouse/_mocks/gateways"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stripe/stripe-go"
)

func TestInsertCustomerSuccess(t *testing.T) {
	assert := assert.New(t)
	mocks, customer := getMockCustomer()

	user := getMockUser()
	req := getMockCustomerRequest(user.ID)

	mocks.sql.ExpectQuery("SELECT userId, stripeCustomerId FROM customer_account").
		WithArgs(user.ID.String()).
		WillReturnRows(getCustomerRows())
	mocks.tc.On("GetUser", user.ID).Return(user, nil)
	mocks.stripe.On("NewCustomer", req.Token, req.UserID.String()).Return("customerID", nil)
	mocks.sql.ExpectPrepare("INSERT INTO customer_account").
		ExpectExec().
		WillReturnResult(sqlmock.NewResult(1, 1))

	c, err := customer.Insert(req)

	assert.NoError(err)
	assert.NotNil(c)
}

func TestInsertExistingSuccess(t *testing.T) {
	assert := assert.New(t)
	mocks, customer := getMockCustomer()

	user := getMockUser()
	req := getMockCustomerRequest(user.ID)
	cust := &models.Customer{
		CustomerID: "customerID",
		UserID:     user.ID,
	}
	c := &stripe.Customer{}

	mocks.sql.ExpectQuery("SELECT userId, stripeCustomerId FROM customer_account").
		WithArgs(user.ID.String()).
		WillReturnRows(getCustomerRows().
			AddRow(user.ID.String(), cust.CustomerID))
	mocks.stripe.On("GetCustomer", cust.CustomerID).Return(c, nil)
	mocks.stripe.On("AddSource", cust.CustomerID, req.Token).Return(c, nil)

	_, err := customer.Insert(req)

	assert.NoError(err)
}

func TestInsertCustomerGetError(t *testing.T) {
	assert := assert.New(t)
	mocks, customer := getMockCustomer()

	user := getMockUser()
	req := getMockCustomerRequest(user.ID)

	mocks.sql.ExpectQuery("SELECT userId, stripeCustomerId FROM customer_account").
		WithArgs(user.ID.String()).
		WillReturnError(fmt.Errorf("some error"))

	_, err := customer.Insert(req)

	assert.Error(err)
}

func TestInsertCustomerUserError(t *testing.T) {
	assert := assert.New(t)
	mocks, customer := getMockCustomer()

	user := getMockUser()
	req := getMockCustomerRequest(user.ID)

	mocks.sql.ExpectQuery("SELECT userId, stripeCustomerId FROM customer_account").
		WithArgs(user.ID.String()).
		WillReturnRows(getCustomerRows())
	mocks.tc.On("GetUser", user.ID).Return(nil, fmt.Errorf("some error"))

	c, err := customer.Insert(req)

	assert.Error(err)
	assert.Nil(c)
}

func TestInsertCustomerUserNil(t *testing.T) {
	assert := assert.New(t)
	mocks, customer := getMockCustomer()

	user := getMockUser()
	req := getMockCustomerRequest(user.ID)

	mocks.sql.ExpectQuery("SELECT userId, stripeCustomerId FROM customer_account").
		WithArgs(user.ID.String()).
		WillReturnRows(getCustomerRows())
	mocks.tc.On("GetUser", user.ID).Return(nil, nil)

	c, err := customer.Insert(req)

	assert.Error(err)
	assert.Nil(c)
}

func TestInsertCustomerStripeError(t *testing.T) {
	assert := assert.New(t)
	mocks, customer := getMockCustomer()

	user := getMockUser()
	req := getMockCustomerRequest(user.ID)

	mocks.sql.ExpectQuery("SELECT userId, stripeCustomerId FROM customer_account").
		WithArgs(user.ID.String()).
		WillReturnRows(getCustomerRows())
	mocks.tc.On("GetUser", user.ID).Return(user, nil)
	mocks.stripe.On("NewCustomer", req.Token, req.UserID.String()).Return("", fmt.Errorf("some error"))

	c, err := customer.Insert(req)

	assert.Error(err)
	assert.Nil(c)
}

func TestInsertCustomerError(t *testing.T) {
	assert := assert.New(t)
	mocks, customer := getMockCustomer()

	user := getMockUser()
	req := getMockCustomerRequest(user.ID)

	mocks.sql.ExpectQuery("SELECT userId, stripeCustomerId FROM customer_account").
		WithArgs(user.ID.String()).
		WillReturnRows(getCustomerRows())
	mocks.tc.On("GetUser", user.ID).Return(user, nil)
	mocks.stripe.On("NewCustomer", req.Token, req.UserID.String()).Return("customerID", nil)
	mocks.sql.ExpectPrepare("INSERT INTO customer_account").
		ExpectExec().
		WillReturnError(fmt.Errorf("some error"))

	c, err := customer.Insert(req)

	assert.Error(err)
	assert.Nil(c)
}

func TestGetCustomerSuccess(t *testing.T) {
	assert := assert.New(t)
	mocks, customer := getMockCustomer()

	user := getMockUser()
	c := &stripe.Customer{}

	mocks.sql.ExpectQuery("SELECT userId, stripeCustomerId FROM customer_account").
		WithArgs(user.ID.String()).
		WillReturnRows(getCustomerRows().
			AddRow(user.ID.String(), "customerID"))
	mocks.stripe.On("GetCustomer", "customerID").Return(c, nil)

	res, err := customer.Get(user.ID)

	assert.NoError(err)
	assert.NotNil(res)
}

func TestGetCustomerError(t *testing.T) {
	assert := assert.New(t)
	mocks, customer := getMockCustomer()

	user := getMockUser()

	mocks.sql.ExpectQuery("SELECT userId, stripeCustomerId FROM customer_account").
		WithArgs(user.ID.String()).
		WillReturnError(fmt.Errorf("some error"))

	res, err := customer.Get(user.ID)

	assert.Error(err)
	assert.Nil(res)
}

func TestSubscribeSuccess(t *testing.T) {
	assert := assert.New(t)
	mocks, customer := getMockCustomer()

	user := getMockUser()
	c := &stripe.Customer{}
	roaster := &models.Roaster{
		Secret: "test",
	}
	plan := &models.Plan{
		PlanIDs: []string{uuid.New()},
	}
	freq := models.Frequency(models.WEEKLY)

	mocks.sql.ExpectQuery("SELECT userId, stripeCustomerId FROM customer_account").
		WithArgs(user.ID.String()).
		WillReturnRows(getCustomerRows().
			AddRow(user.ID.String(), "customerID"))
	mocks.sql.ExpectPrepare("INSERT INTO subscribed").
		ExpectExec().
		WithArgs("customerID", "connectedID", roaster.ID.String(), "sub").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mocks.stripe.On("Subscribe", roaster, "customerID", plan.PlanIDs[0], uint64(1)).Return("connectedID", &stripe.Sub{ID: "sub"}, nil)
	mocks.stripe.On("GetCustomer", "customerID").Return(c, nil)

	err := customer.Subscribe(user.ID, roaster, plan, freq, uint64(1))

	assert.NoError(err)
}

func TestSubscribeUserFail(t *testing.T) {
	assert := assert.New(t)
	mocks, customer := getMockCustomer()

	user := getMockUser()
	roaster := &models.Roaster{
		Secret: "test",
	}
	plan := &models.Plan{
		PlanIDs: []string{uuid.New()},
	}
	freq := models.Frequency(models.WEEKLY)

	mocks.sql.ExpectQuery("SELECT userId, stripeCustomerId FROM customer_account").
		WithArgs(user.ID.String()).
		WillReturnError(fmt.Errorf("some error"))
	err := customer.Subscribe(user.ID, roaster, plan, freq, uint64(1))

	assert.Error(err)
}

func TestSubscribeFreqFail(t *testing.T) {
	assert := assert.New(t)
	mocks, customer := getMockCustomer()

	user := getMockUser()
	c := &stripe.Customer{}
	plan := &models.Plan{
		PlanIDs: []string{uuid.New()},
	}
	roaster := &models.Roaster{
		Secret: "test",
	}
	freq := models.Frequency("badstring")

	mocks.sql.ExpectQuery("SELECT userId, stripeCustomerId FROM customer_account").
		WithArgs(user.ID.String()).
		WillReturnRows(getCustomerRows().
			AddRow(user.ID.String(), "customerID"))
	mocks.stripe.On("GetCustomer", "customerID").Return(c, nil)

	err := customer.Subscribe(user.ID, roaster, plan, freq, uint64(1))

	assert.Error(err)
}

func TestSubscribeFail(t *testing.T) {
	assert := assert.New(t)
	mocks, customer := getMockCustomer()

	user := getMockUser()
	c := &stripe.Customer{}
	plan := &models.Plan{
		PlanIDs: []string{uuid.New(), uuid.New()},
	}
	roaster := &models.Roaster{
		Secret: "test",
	}
	freq := models.Frequency(models.BIWEEKLY)

	mocks.sql.ExpectQuery("SELECT userId, stripeCustomerId FROM customer_account").
		WithArgs(user.ID.String()).
		WillReturnRows(getCustomerRows().
			AddRow(user.ID.String(), "customerID"))
	mocks.stripe.On("Subscribe", roaster, "customerID", plan.PlanIDs[1], uint64(1)).Return("", &stripe.Sub{ID: "sub"}, fmt.Errorf("some error"))
	mocks.stripe.On("GetCustomer", "customerID").Return(c, nil)

	err := customer.Subscribe(user.ID, roaster, plan, freq, uint64(1))

	assert.Error(err)
}

func TestDeleteSuccess(t *testing.T) {
	assert := assert.New(t)
	mocks, customer := getMockCustomer()

	user := getMockUser()
	c := &stripe.Customer{}

	mocks.sql.ExpectQuery("SELECT userId, stripeCustomerId FROM customer_account").
		WithArgs(user.ID.String()).
		WillReturnRows(getCustomerRows().
			AddRow(user.ID.String(), "customerID"))
	mocks.sql.ExpectPrepare("DELETE FROM customer_account").
		ExpectExec().
		WithArgs(user.ID.String()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mocks.stripe.On("GetCustomer", "customerID").Return(c, nil)
	mocks.stripe.On("DeleteCustomer", "customerID").Return(nil)

	err := customer.Delete(user.ID)

	assert.NoError(err)
}

func TestDeleteUserFail(t *testing.T) {
	assert := assert.New(t)
	mocks, customer := getMockCustomer()

	user := getMockUser()

	mocks.sql.ExpectQuery("SELECT userId, stripeCustomerId FROM customer_account").
		WithArgs(user.ID.String()).
		WillReturnError(fmt.Errorf("some error"))

	err := customer.Delete(user.ID)

	assert.Error(err)
}

func TestDeleteStripeFail(t *testing.T) {
	assert := assert.New(t)
	mocks, customer := getMockCustomer()

	user := getMockUser()
	c := &stripe.Customer{}

	mocks.sql.ExpectQuery("SELECT userId, stripeCustomerId FROM customer_account").
		WithArgs(user.ID.String()).
		WillReturnRows(getCustomerRows().
			AddRow(user.ID.String(), "customerID"))
	mocks.stripe.On("GetCustomer", "customerID").Return(c, nil)
	mocks.stripe.On("DeleteCustomer", "customerID").Return(fmt.Errorf("some error"))

	err := customer.Delete(user.ID)

	assert.Error(err)
}

type mockContext struct {
	sql       sqlmock.Sqlmock
	stripe    *mocks.Stripe
	tc        *tmocks.TownCenterI
	warehouse *wmocks.Warehouse
}

func getMockCustomer() (*mockContext, *Customer) {
	s, mock, _ := sqlmock.New()
	mocks := &mockContext{
		sql:    mock,
		stripe: &mocks.Stripe{},
		tc:     &tmocks.TownCenterI{},
	}
	return mocks, NewCustomer(&bgateways.MySQL{DB: s}, mocks.stripe, mocks.tc)
}

func getMockCustomerRequest(id uuid.UUID) *models.CustomerRequest {
	return &models.CustomerRequest{
		Token:  "token",
		UserID: id,
	}
}

func getMockUser() *tmodels.User {
	return tmodels.NewUser(
		"passHash", "firstName",
		"lastName", "email",
		"phone", "addressLine1",
		"addressLine2", "addressCity",
		"addressState", "addressZip",
		"addressCountry")
}

func getCustomerRows() sqlmock.Rows {
	return sqlmock.NewRows([]string{"userId", "stripeCustomerId"})
}
