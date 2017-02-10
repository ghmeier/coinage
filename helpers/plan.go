package helpers

import (
	"strings"

	"github.com/pborman/uuid"
	"github.com/stripe/stripe-go"

	g "github.com/ghmeier/bloodlines/gateways"
	"github.com/jonnykry/coinage/gateways"
	"github.com/jonnykry/coinage/models"
	item "github.com/lcollin/warehouse/models"
)

type Plan struct {
	*baseHelper
	Stripe gateways.Stripe
}

func NewPlan(sql g.SQL, stripe gateways.Stripe) *Plan {
	return &Plan{
		baseHelper: &baseHelper{sql: sql},
		Stripe:     stripe,
	}
}

func (p *Plan) Insert(id uuid.UUID, accountID string, req *models.PlanRequest) (*models.Plan, error) {

	// TODO: get item from warehouse by id
	item := &item.Item{}

	planIDs := make([]string, 0)
	plans := make([]*stripe.Plan, 0)

	for i := 0; i < len(models.Frequencies); i++ {
		stripe, err := p.Stripe.NewPlan(accountID, item, models.Frequencies[i])
		if err != nil {
			return nil, err
		}

		plans = append(plans, stripe)
		planIDs = append(planIDs, stripe.ID)
	}

	plan := models.NewPlan(id, req.ItemID, planIDs)
	plan.Plans = plans
	err := p.sql.Modify("INSERT INTO plan (roasterId,itemId,planIds)VALUES(?,?)",
		plan.RoasterID,
		plan.ItemID,
		strings.Join(plan.PlanIDs, ","),
	)
	if err != nil {
		return nil, err
	}
	return plan, nil
}

func (p *Plan) GetByRoaster(id uuid.UUID, offset int, limit int) ([]*models.Plan, error) {
	rows, err := p.sql.Select("SELECT roasterId,itemId,planIds FROM plan WHERE roasterId=? ORDER BY itemId ASC LIMIT ?,?",
		id,
		offset,
		limit,
	)
	if err != nil {
		return nil, err
	}

	plans, _ := models.PlanFromSQL(rows)

	return plans, nil
}

func (p *Plan) Get(id uuid.UUID, itemID uuid.UUID) (*models.Plan, error) {
	rows, err := p.sql.Select("SELECT roasterId,itemId,planIds FROM plan WHERE roasterId=? AND itemId=?",
		id,
		itemID,
	)

	if err != nil {
		return nil, err
	}

	plans, _ := models.PlanFromSQL(rows)

	return plans[0], nil
}

func (p *Plan) Update(id string, itemId uuid.UUID) (*models.Plan, error) {
	// err := p.sql.Modify("UPDATE plan SET roasterId=?,itemId=?,planIds=? WHERE itemId=?",
	// 	plan.RoasterID,
	// 	plan.ItemID,
	// 	strings.Join(plan.PlanIDs, ","),
	// )
	return nil, nil
}

func (p *Plan) Delete(id, planID string) error {
	return nil
}
