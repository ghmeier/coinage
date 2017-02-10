package helpers

import (
	"fmt"
	"testing"

	bgateways "github.com/ghmeier/bloodlines/gateways"
	tmocks "github.com/jakelong95/TownCenter/_mocks"
	mocks "github.com/jonnykry/coinage/_mocks/gateways"
	"github.com/jonnykry/coinage/models"
	item "github.com/lcollin/warehouse/models"
	cmocks "github.com/yuderekyu/covenant/_mocks/gateways"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stripe/stripe-go"
)

func TestInsertPlanSuccess(t *testing.T) {
	assert := assert.New(t)
	mocks, plan := getMockPHelper()

	id := uuid.NewUUID()
	accountID := "accountID"
	plans := getMockPlans()
	req := getMockPlanRequest()

	mocks.stripe.On("NewPlan", accountID, &item.Item{}, models.WEEKLY).Return(plans[0], nil)
	mocks.stripe.On("NewPlan", accountID, &item.Item{}, models.BIWEEKLY).Return(plans[1], nil)
	mocks.stripe.On("NewPlan", accountID, &item.Item{}, models.TRIWEEKLY).Return(plans[2], nil)
	mocks.stripe.On("NewPlan", accountID, &item.Item{}, models.MONTHLY).Return(plans[3], nil)
	mocks.sql.ExpectPrepare("INSERT INTO plan").
		ExpectExec().
		WillReturnResult(sqlmock.NewResult(1, 1))

	c, err := plan.Insert(id, accountID, req)

	assert.NoError(err)
	assert.NotNil(c)
}

func TestInsertPlanSQLFail(t *testing.T) {
	assert := assert.New(t)
	mocks, plan := getMockPHelper()

	id := uuid.NewUUID()
	accountID := "accountID"
	plans := getMockPlans()
	req := getMockPlanRequest()

	mocks.stripe.On("NewPlan", accountID, &item.Item{}, models.WEEKLY).Return(plans[0], nil)
	mocks.stripe.On("NewPlan", accountID, &item.Item{}, models.BIWEEKLY).Return(plans[1], nil)
	mocks.stripe.On("NewPlan", accountID, &item.Item{}, models.TRIWEEKLY).Return(plans[2], nil)
	mocks.stripe.On("NewPlan", accountID, &item.Item{}, models.MONTHLY).Return(plans[3], nil)
	mocks.sql.ExpectPrepare("INSERT INTO plan").
		ExpectExec().
		WillReturnError(fmt.Errorf("some error"))

	c, err := plan.Insert(id, accountID, req)

	assert.Error(err)
	assert.Nil(c)
}

func getMockPHelper() (*mockContext, *Plan) {
	s, mock, _ := sqlmock.New()
	mocks := &mockContext{
		sql:    mock,
		stripe: &mocks.Stripe{},
		tc:     &tmocks.TownCenterI{},
		c:      &cmocks.Covenant{},
	}
	return mocks, NewPlan(&bgateways.MySQL{DB: s}, mocks.stripe)
}

func getMockPlan(id uuid.UUID) *models.Plan {
	return &models.Plan{
		RoasterID: uuid.NewUUID(),
		ItemID:    uuid.NewUUID(),
		PlanIDs:   []string{"1", "2", "3", "4"},
	}
}

func getMockPlanRequest() *models.PlanRequest {
	return &models.PlanRequest{
		ItemID: uuid.NewUUID(),
	}
}

func getMockPlans() []*stripe.Plan {
	return []*stripe.Plan{
		&stripe.Plan{ID: "1"},
		&stripe.Plan{ID: "2"},
		&stripe.Plan{ID: "3"},
		&stripe.Plan{ID: "4"},
	}
}

func getPlanRows() sqlmock.Rows {
	return sqlmock.NewRows([]string{"roasterId", "itemId", "planIds"})
}
