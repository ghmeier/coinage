package helpers

import (
	"fmt"
	"strings"
	"testing"

	bgateways "github.com/ghmeier/bloodlines/gateways"
	mocks "github.com/ghmeier/coinage/_mocks/gateways"
	"github.com/ghmeier/coinage/models"
	tmocks "github.com/jakelong95/TownCenter/_mocks"
	item "github.com/lcollin/warehouse/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stripe/stripe-go"
)

func TestInsertPlanSuccess(t *testing.T) {
	assert := assert.New(t)
	mocks, p := getMockPHelper()

	roast := getMockCRoaster()
	plans := getMockPlans()
	req := getMockPlanRequest()

	mocks.stripe.On("NewPlan", roast.Secret, &item.Item{}, models.Frequencies[0]).Return(plans[0], nil)
	mocks.stripe.On("NewPlan", roast.Secret, &item.Item{}, models.Frequencies[1]).Return(plans[1], nil)
	mocks.stripe.On("NewPlan", roast.Secret, &item.Item{}, models.Frequencies[2]).Return(plans[2], nil)
	mocks.stripe.On("NewPlan", roast.Secret, &item.Item{}, models.Frequencies[3]).Return(plans[3], nil)
	mocks.sql.ExpectPrepare("INSERT INTO plan").
		ExpectExec().
		WillReturnResult(sqlmock.NewResult(1, 1))

	c, err := p.Insert(roast, req)

	assert.NoError(err)
	assert.NotNil(c)
}

func TestInsertPlanSQLFail(t *testing.T) {
	assert := assert.New(t)
	mocks, p := getMockPHelper()

	roast := getMockCRoaster()
	plans := getMockPlans()
	req := getMockPlanRequest()

	mocks.stripe.On("NewPlan", roast.Secret, &item.Item{}, models.Frequencies[0]).Return(plans[0], nil)
	mocks.stripe.On("NewPlan", roast.Secret, &item.Item{}, models.Frequencies[1]).Return(plans[1], nil)
	mocks.stripe.On("NewPlan", roast.Secret, &item.Item{}, models.Frequencies[2]).Return(plans[2], nil)
	mocks.stripe.On("NewPlan", roast.Secret, &item.Item{}, models.Frequencies[3]).Return(plans[3], nil)
	mocks.sql.ExpectPrepare("INSERT INTO plan").
		ExpectExec().
		WillReturnError(fmt.Errorf("some error"))

	c, err := p.Insert(roast, req)

	assert.Error(err)
	assert.Nil(c)
}

func TestInsertPlanStripeFail(t *testing.T) {
	assert := assert.New(t)
	mocks, p := getMockPHelper()

	roast := getMockCRoaster()
	req := getMockPlanRequest()

	mocks.stripe.On("NewPlan", roast.Secret, &item.Item{}, models.Frequencies[0]).Return(nil, fmt.Errorf("some error"))

	c, err := p.Insert(roast, req)

	assert.Error(err)
	assert.Nil(c)
}

func TestGetPlanSuccess(t *testing.T) {
	assert := assert.New(t)
	mocks, p := getMockPHelper()

	roast := getMockCRoaster()
	plans := getMockPlans()
	plan := getMockPlan(roast.ID)

	req := getMockPlanRequest()

	mocks.stripe.On("GetPlan", roast.Secret, plan.PlanIDs[0]).Return(plans[0], nil)
	mocks.stripe.On("GetPlan", roast.Secret, plan.PlanIDs[1]).Return(plans[1], nil)
	mocks.stripe.On("GetPlan", roast.Secret, plan.PlanIDs[2]).Return(plans[2], nil)
	mocks.stripe.On("GetPlan", roast.Secret, plan.PlanIDs[3]).Return(plans[3], nil)
	mocks.sql.ExpectQuery("SELECT roasterId,itemId,planIds FROM plan").
		WillReturnRows(getPlanRows().
			AddRow(roast.ID.String(), req.ItemID.String(), strings.Join(plan.PlanIDs, ",")))

	c, err := p.Get(roast, req.ItemID)

	assert.NoError(err)
	assert.NotNil(c)
}

func TestGetPlanSQLFail(t *testing.T) {
	assert := assert.New(t)
	mocks, p := getMockPHelper()

	roaster := getMockCRoaster()
	req := getMockPlanRequest()

	mocks.sql.ExpectQuery("SELECT roasterId,itemId,planIds FROM plan").
		WillReturnError(fmt.Errorf("some error"))

	c, err := p.Get(roaster, req.ItemID)

	assert.Error(err)
	assert.Nil(c)
}

func TestGetPlanStripeFail(t *testing.T) {
	assert := assert.New(t)
	mocks, p := getMockPHelper()

	roast := getMockCRoaster()
	plan := getMockPlan(roast.ID)

	req := getMockPlanRequest()

	mocks.stripe.On("GetPlan", roast.Secret, plan.PlanIDs[0]).Return(nil, fmt.Errorf("some error"))
	mocks.stripe.On("GetPlan", roast.Secret, plan.PlanIDs[1]).Return(nil, fmt.Errorf("some error"))
	mocks.stripe.On("GetPlan", roast.Secret, plan.PlanIDs[2]).Return(nil, fmt.Errorf("some error"))
	mocks.stripe.On("GetPlan", roast.Secret, plan.PlanIDs[3]).Return(nil, fmt.Errorf("some error"))
	mocks.sql.ExpectQuery("SELECT roasterId,itemId,planIds FROM plan").
		WillReturnRows(getPlanRows().
			AddRow(roast.ID.String(), req.ItemID.String(), strings.Join(plan.PlanIDs, ",")))

	c, err := p.Get(roast, req.ItemID)

	assert.Error(err)
	assert.Nil(c)
}

func TestGetPlanByRoasterSuccess(t *testing.T) {
	assert := assert.New(t)
	mocks, p := getMockPHelper()

	roaster := getMockCRoaster()
	plans := getMockPlans()
	plan := getMockPlan(roaster.ID)

	mocks.stripe.On("GetPlan", roaster.Secret, plan.PlanIDs[0]).Return(plans[0], nil)
	mocks.stripe.On("GetPlan", roaster.Secret, plan.PlanIDs[1]).Return(plans[1], nil)
	mocks.stripe.On("GetPlan", roaster.Secret, plan.PlanIDs[2]).Return(plans[2], nil)
	mocks.stripe.On("GetPlan", roaster.Secret, plan.PlanIDs[3]).Return(plans[3], nil)
	// mocks.stripe.On("GetPlan", roaster.AccountID, plan.PlanIDs[0]).Return(plans[0], nil)
	// mocks.stripe.On("GetPlan", roaster.AccountID, plan.PlanIDs[1]).Return(plans[1], nil)
	// mocks.stripe.On("GetPlan", roaster.AccountID, plan.PlanIDs[2]).Return(plans[2], nil)
	// mocks.stripe.On("GetPlan", roaster.AccountID, plan.PlanIDs[3]).Return(plans[3], nil)
	mocks.sql.ExpectQuery("SELECT roasterId,itemId,planIds FROM plan").
		WillReturnRows(getPlanRows().
			AddRow(roaster.ID.String(), uuid.New(), strings.Join(plan.PlanIDs, ",")).
			AddRow(roaster.ID.String(), uuid.New(), strings.Join(plan.PlanIDs, ",")))

	c, err := p.GetByRoaster(roaster, 0, 20)

	assert.NoError(err)
	assert.NotNil(c)
	assert.Equal(2, len(c))
}

func TestGetPlanByRoasterFail(t *testing.T) {
	assert := assert.New(t)
	mocks, p := getMockPHelper()

	roaster := getMockCRoaster()

	mocks.sql.ExpectQuery("SELECT roasterId,itemId,planIds FROM plan").
		WillReturnError(fmt.Errorf("some error"))
	c, err := p.GetByRoaster(roaster, 0, 20)

	assert.Error(err)
	assert.Nil(c)
}

func getMockPHelper() (*mockContext, *Plan) {
	s, mock, _ := sqlmock.New()
	mocks := &mockContext{
		sql:    mock,
		stripe: &mocks.Stripe{},
		tc:     &tmocks.TownCenterI{},
	}
	return mocks, NewPlan(&bgateways.MySQL{DB: s}, mocks.stripe)
}

func getMockPlan(id uuid.UUID) *models.Plan {
	return &models.Plan{
		RoasterID: id,
		ItemID:    uuid.NewUUID(),
		PlanIDs:   []string{"1", "2", "3", "4"},
	}
}

func getMockCRoaster() *models.Roaster {
	return &models.Roaster{
		ID:          uuid.NewUUID(),
		AccountID:   "test",
		Secret:      "secret",
		Publishable: "test",
	}
}

func getMockPlanRequest() *models.PlanRequest {
	return &models.PlanRequest{
		ItemID: uuid.NewUUID(),
	}
}

func getMockPlans() []*stripe.Plan {
	return []*stripe.Plan{
		{ID: "1"},
		{ID: "2"},
		{ID: "3"},
		{ID: "4"},
	}
}

func getPlanRows() sqlmock.Rows {
	return sqlmock.NewRows([]string{"roasterId", "itemId", "planIds"})
}
