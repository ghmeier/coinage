package helpers

import (
	"database/sql"
	"strings"

	"github.com/pborman/uuid"
	"github.com/stripe/stripe-go"

	g "github.com/ghmeier/bloodlines/gateways"
	"github.com/ghmeier/coinage/gateways"
	"github.com/ghmeier/coinage/models"
	item "github.com/lcollin/warehouse/models"
)

/*Plan manages retrieval and manipulating roaster plan information*/
type Plan struct {
	*base
	Stripe gateways.Stripe
}

/*NewPlan initializes and returns a plan with the given gateways*/
func NewPlan(sql g.SQL, stripe gateways.Stripe) *Plan {
	return &Plan{
		base:   &base{sql: sql},
		Stripe: stripe,
	}
}

/*Insert creates a new roaster plan in stripe and adds a record to the db*/
func (p *Plan) Insert(roaster *models.Roaster, req *models.PlanRequest) (*models.Plan, error) {

	// TODO: get item from warehouse by id
	item := &item.Item{}

	planIDs := make([]string, 0)
	plans := make([]*stripe.Plan, 0)

	for i := 0; i < len(models.Frequencies); i++ {
		stripe, err := p.Stripe.NewPlan(roaster.Secret, item, models.Frequencies[i])
		if err != nil {
			return nil, err
		}

		plans = append(plans, stripe)
		planIDs = append(planIDs, stripe.ID)
	}

	plan := models.NewPlan(roaster.ID, req.ItemID, planIDs)
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

/*GetByRoaster returns all plans associated with a roaster*/
func (p *Plan) GetByRoaster(roaster *models.Roaster, offset int, limit int) ([]*models.Plan, error) {
	rows, err := p.sql.Select("SELECT roasterId,itemId,planIds FROM plan WHERE roasterId=? ORDER BY itemId ASC LIMIT ?,?",
		roaster.ID,
		offset,
		limit,
	)
	if err != nil {
		return nil, err
	}

	return p.plan(roaster.Secret, rows)
}

/*Get returns a plan associated with a roaster and itemid*/
func (p *Plan) Get(roaster *models.Roaster, itemID uuid.UUID) (*models.Plan, error) {
	rows, err := p.sql.Select("SELECT roasterId,itemId,planIds FROM plan WHERE roasterId=? AND itemId=?",
		roaster.ID,
		itemID,
	)
	if err != nil {
		return nil, err
	}

	plans, err := p.plan(roaster.Secret, rows)
	if err != nil {
		return nil, err
	}

	return plans[0], err
}

func (p *Plan) plan(secret string, rows *sql.Rows) ([]*models.Plan, error) {
	plans, _ := models.PlanFromSQL(rows)

	for i := range plans {
		stripePlans, err := p.plans(secret, plans[i].PlanIDs)
		if err != nil {
			return nil, err
		}
		plans[i].Plans = stripePlans
	}

	return plans, nil
}

func (p *Plan) plans(secret string, ids []string) ([]*stripe.Plan, error) {
	plans := make([]*stripe.Plan, 0)
	var err error
	var plan *stripe.Plan
	for i := 0; i < len(models.Frequencies); i++ {
		plan, err = p.Stripe.GetPlan(secret, ids[i])
		plans = append(plans, plan)
	}
	return plans, err
}

/*Update is not implemented*/
func (p *Plan) Update(id string, itemID uuid.UUID) (*models.Plan, error) {
	// err := p.sql.Modify("UPDATE plan SET roasterId=?,itemId=?,planIds=? WHERE itemId=?",
	// 	plan.RoasterID,
	// 	plan.ItemID,
	// 	strings.Join(plan.PlanIDs, ","),
	// )
	return nil, nil
}

/*Delete is not implemented*/
func (p *Plan) Delete(id, planID string) error {
	return nil
}
