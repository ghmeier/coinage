package models

import (
	"database/sql"

	"github.com/stripe/stripe-go"
)

type Plan struct {
	RoasterID string       `json:"roasterId"`
	PlanID    string       `json:"planId"`
	Plan      *stripe.Plan `json:"plan"`
}

type PlanRequest struct {
}

func NewPlan(roasterID, planID string) *Plan {
	return &Plan{
		RoasterID: roasterID,
		PlanID:    planID,
	}
}

func PlanFromSQL(rows *sql.Rows) ([]*Plan, error) {
	plans := make([]*Plan, 0)

	for rows.Next() {
		p := &Plan{}
		rows.Scan(&p.RoasterID, &p.PlanID)
		plans = append(plans, p)
	}

	return plans, nil
}
