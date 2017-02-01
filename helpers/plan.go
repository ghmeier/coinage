package helpers

import (
	"github.com/pborman/uuid"

	"github.com/ghmeier/bloodlines/gateways"
	"github.com/jonnykry/coinage/models"
)

type Plan struct {
	*baseHelper
	Stripe gateways.Stripe
}

func NewPlan(sql gateways.SQL, stripe gateways.Stripe) *Plan {
	return &Plan{
		baseHelper: &baseHelper{sql: sql},
		Stripe:     stripe,
	}
}

func (p *Plan) Insert(plan *models.Plan) error {
	err := p.sql.Modify("INSERT INTO plan (roasterId,planId)VALUES(?,?)",
		subscription.RoasterID,
		subscription.PlanID,
	)
	return err
}

func (p *Plan) GetByRoaster(id uuid.UUID, offset int, limit, int) ([]*models.Plan, error) {
	rows, err := p.Select("SELECT roasterId,planId FROM plan WHERE roasterId=? ORDER BY planId ASC LIMIT ?,?",
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

func (p *Plan) GetPlan(id, planID uuid.UUID) (*models.Plan, error) {
	return nil, nil
}

func (b *Plan) Update(plan *models.Plan) error {
	err := b.sql.Modify("UPDATE plan SET roasterId=?,planId=? WHERE roasterId=?",
		plan.RoasterID,
		plan.PlanID,
	)
	return err
}
