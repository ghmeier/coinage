package helpers

import (
	"github.com/pborman/uuid"

	g "github.com/ghmeier/bloodlines/gateways"
	"github.com/jonnykry/coinage/gateways"
	"github.com/jonnykry/coinage/models"
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

func (p *Plan) Insert(id string, req *models.PlanRequest) (*models.Plan, error) {
	stripe, err := p.Stripe.NewPlan(id, req)
	if err != nil {
		return nil, err
	}

	plan := models.NewPlan(id, stripe.ID)
	plan.Plan = stripe
	err = p.sql.Modify("INSERT INTO plan (roasterId,planId)VALUES(?,?)",
		plan.RoasterID,
		plan.PlanID,
	)
	return plan, err
}

func (p *Plan) GetByRoaster(id uuid.UUID, offset int, limit int) ([]*models.Plan, error) {
	rows, err := p.sql.Select("SELECT roasterId,planId FROM plan WHERE roasterId=? ORDER BY planId ASC LIMIT ?,?",
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

func (p *Plan) Get(id string, planID string) (*models.Plan, error) {
	stripe, err := p.Stripe.GetPlan(id, planID)
	if err != nil {
		return nil, err
	}

	plan := models.NewPlan(id, planID)
	plan.Plan = stripe
	return plan, nil
}

func (p *Plan) Update(plan *models.Plan, req *models.PlanRequest) (*models.Plan, error) {
	err := p.sql.Modify("UPDATE plan SET roasterId=?,planId=? WHERE roasterId=?",
		plan.RoasterID,
		plan.PlanID,
	)
	return plan, err
}

func (p *Plan) Delete(id, planID string) error {
	return nil
}
