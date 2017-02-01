package models

import (
	"database/sql"
)

type Plan struct {
	RoasterID string `json:"roasterId"`
	PlanID    string `json:"planId"`
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
		rows.Scan(&p.RoasterID, &c.PlanID)
		roasterAccount = append(plans, p)
	}

	return plans, nil
}
