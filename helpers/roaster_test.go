package helpers

import (
	"fmt"
	"testing"

	bgateways "github.com/ghmeier/bloodlines/gateways"
	mocks "github.com/ghmeier/coinage/_mocks/gateways"
	"github.com/ghmeier/coinage/models"
	tmocks "github.com/jakelong95/TownCenter/_mocks"
	tmodels "github.com/jakelong95/TownCenter/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stripe/stripe-go"
)

func TestInsertRoasterSuccess(t *testing.T) {
	assert := assert.New(t)
	mocks, roaster := getMockRHelper()

	user := getMockRUser()
	r := getMockRoaster(user.RoasterId)
	req := getMockRoasterRequest(user.ID)

	mocks.tc.On("GetUser", user.ID).Return(user, nil)
	mocks.tc.On("GetRoaster", r.ID).Return(r, nil)
	mocks.stripe.On("NewAccount", req.Country, user, r).Return(&stripe.Account{ID: "stripeID"}, nil)
	mocks.sql.ExpectPrepare("INSERT INTO roaster_account").
		ExpectExec().
		WillReturnResult(sqlmock.NewResult(1, 1))

	c, err := roaster.Insert(req)

	assert.NoError(err)
	assert.NotNil(c)
}

func TestInsertRoasterUserFail(t *testing.T) {
	assert := assert.New(t)
	mocks, roaster := getMockRHelper()

	user := getMockRUser()
	req := getMockRoasterRequest(user.ID)

	mocks.tc.On("GetUser", user.ID).Return(nil, fmt.Errorf("some error"))
	c, err := roaster.Insert(req)
	assert.Error(err)
	assert.Nil(c)

	mocks.tc.On("GetUser", user.ID).Return(nil, nil)
	c, err = roaster.Insert(req)
	assert.Error(err)
	assert.Nil(c)

	user.RoasterId = nil
	mocks.tc.On("GetUser", user.ID).Return(user, nil)
	c, err = roaster.Insert(req)
	assert.Error(err)
	assert.Nil(c)
}

func TestInsertRoasterRoasterFail(t *testing.T) {
	assert := assert.New(t)
	mocks, roaster := getMockRHelper()

	user := getMockRUser()
	req := getMockRoasterRequest(user.ID)

	mocks.tc.On("GetUser", user.ID).Return(user, nil)
	mocks.tc.On("GetRoaster", user.RoasterId).Return(nil, fmt.Errorf("some error"))

	c, err := roaster.Insert(req)

	assert.Error(err)
	assert.Nil(c)

	mocks.tc.On("GetUser", user.ID).Return(user, nil)
	mocks.tc.On("GetRoaster", user.RoasterId).Return(nil, nil)

	c, err = roaster.Insert(req)

	assert.Error(err)
	assert.Nil(c)
}

func TestInsertRoasterStripeFail(t *testing.T) {
	assert := assert.New(t)
	mocks, roaster := getMockRHelper()

	user := getMockRUser()
	r := getMockRoaster(user.RoasterId)
	req := getMockRoasterRequest(user.ID)

	mocks.tc.On("GetUser", user.ID).Return(user, nil)
	mocks.tc.On("GetRoaster", r.ID).Return(r, nil)
	mocks.stripe.On("NewAccount", req.Country, user, r).Return(nil, fmt.Errorf("some error"))

	c, err := roaster.Insert(req)

	assert.Error(err)
	assert.Nil(c)
}

func TestInsertRoasterSQLFail(t *testing.T) {
	assert := assert.New(t)
	mocks, roaster := getMockRHelper()

	user := getMockRUser()
	r := getMockRoaster(user.RoasterId)
	req := getMockRoasterRequest(user.ID)

	mocks.tc.On("GetUser", user.ID).Return(user, nil)
	mocks.tc.On("GetRoaster", r.ID).Return(r, nil)
	mocks.stripe.On("NewAccount", req.Country, user, r).Return(&stripe.Account{ID: "stripeID"}, nil)
	mocks.sql.ExpectPrepare("INSERT INTO roater_account").
		ExpectExec().
		WillReturnError(fmt.Errorf("some error"))

	c, err := roaster.Insert(req)

	assert.Error(err)
	assert.Nil(c)
}

func TestGetByUserIDRoasterSuccess(t *testing.T) {
	assert := assert.New(t)
	mocks, roaster := getMockRHelper()

	user := getMockRUser()
	tRoaster := getMockRoaster(user.RoasterId)
	r := getMockRoasterAccount(user.RoasterId)

	mocks.stripe.On("GetAccount", r.AccountID).Return(&stripe.Account{ID: "stripeID"}, nil)
	mocks.tc.On("GetUser", user.ID).Return(user, nil)
	mocks.tc.On("GetRoaster", tRoaster.ID).Return(tRoaster, nil)
	mocks.sql.ExpectQuery("SELECT id, stripeAccountId FROM roaster_account").
		WithArgs(user.RoasterId.String()).
		WillReturnRows(getRoasterRows().
			AddRow(user.RoasterId.String(), r.AccountID))
	c, err := roaster.GetByUserID(user.ID)

	assert.NoError(err)
	assert.NotNil(c)
}

func TestGetByUserIDRoasterFail(t *testing.T) {
	assert := assert.New(t)
	mocks, roaster := getMockRHelper()

	user := getMockRUser()
	tRoaster := getMockRoaster(user.RoasterId)

	mocks.tc.On("GetUser", user.ID).Return(user, nil)
	mocks.tc.On("GetRoaster", tRoaster.ID).Return(tRoaster, nil)
	mocks.sql.ExpectQuery("SELECT id, stripeAccountId FROM roaster_account").
		WithArgs(user.RoasterId.String()).
		WillReturnError(fmt.Errorf("some error"))
	c, err := roaster.GetByUserID(user.ID)

	assert.Error(err)
	assert.Nil(c)
}

func TestGetIDRoasterAccountFail(t *testing.T) {
	assert := assert.New(t)
	mocks, roaster := getMockRHelper()

	id := uuid.NewUUID()
	r := getMockRoasterAccount(id)

	mocks.stripe.On("GetAccount", r.AccountID).Return(nil, fmt.Errorf("some error"))
	mocks.sql.ExpectQuery("SELECT id, stripeAccountId FROM roaster_account").
		WithArgs(id.String()).
		WillReturnRows(getRoasterRows().
			AddRow(r.ID.String(), r.AccountID))
	c, err := roaster.Get(id)

	assert.Error(err)
	assert.Nil(c)
}

func TestGetIDRoasterSQLFail(t *testing.T) {
	assert := assert.New(t)
	mocks, roaster := getMockRHelper()

	id := uuid.NewUUID()

	mocks.sql.ExpectQuery("SELECT id, stripeAccountId FROM roaster_account").
		WithArgs(id.String()).
		WillReturnError(fmt.Errorf("some error"))
	c, err := roaster.Get(id)

	assert.Error(err)
	assert.Nil(c)
}

func getMockRHelper() (*mockContext, *Roaster) {
	s, mock, _ := sqlmock.New()
	mocks := &mockContext{
		sql:    mock,
		stripe: &mocks.Stripe{},
		tc:     &tmocks.TownCenterI{},
	}
	return mocks, NewRoaster(&bgateways.MySQL{DB: s}, mocks.stripe, mocks.tc)
}

func getMockRoasterAccount(id uuid.UUID) *models.Roaster {
	return &models.Roaster{
		ID:        id,
		AccountID: "accountID",
	}
}

func getMockRoasterRequest(id uuid.UUID) *models.RoasterRequest {
	return &models.RoasterRequest{
		Country: "US",
		UserID:  id,
	}
}

func getMockRUser() *tmodels.User {
	u := tmodels.NewUser(
		"passHash", "firstName",
		"lastName", "email",
		"phone", "addressLine1",
		"addressLine2", "addressCity",
		"addressState", "addressZip",
		"addressCountry")
	u.RoasterId = uuid.NewUUID()
	return u
}

func getMockRoaster(id uuid.UUID) *tmodels.Roaster {
	r := tmodels.NewRoaster(
		"name", "email",
		"phone", "addressLine1",
		"addressLine2", "addressCity",
		"addressState", "addressZip",
		"addressCountry")
	r.ID = id
	return r
}

func getRoasterRows() sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "stripeAccountId"})
}
