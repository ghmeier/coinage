package models

import (
	"database/sql"
	"strings"

	"github.com/pborman/uuid"
	"github.com/stripe/stripe-go"
)

type Plan struct {
	RoasterID uuid.UUID      `json:"roasterId"`
	ItemID    uuid.UUID      `json:"itemId"`
	PlanIDs   []string       `json:"planIds"`
	Plans     []*stripe.Plan `json:"plans"`
}

type PlanRequest struct {
	ItemID uuid.UUID `json:"itemId"`
}

func NewPlan(roasterID, itemID uuid.UUID, planIDs []string) *Plan {
	return &Plan{
		RoasterID: roasterID,
		ItemID:    itemID,
		PlanIDs:   planIDs,
	}
}

func PlanFromSQL(rows *sql.Rows) ([]*Plan, error) {
	plans := make([]*Plan, 0)

	for rows.Next() {
		p := &Plan{}
		var planIDs string
		rows.Scan(&p.RoasterID, &p.ItemID, &planIDs)

		p.PlanIDs = strings.Split(planIDs, ",")

		plans = append(plans, p)
	}

	return plans, nil
}

func ToFrequency(s string) (int, bool) {
	switch s {
	case WEEKLY:
		return 0, true
	case BIWEEKLY:
		return 1, true
	case TRIWEEKLY:
		return 2, true
	case MONTHLY:
		return 3, true
	default:
		return -1, false
	}
}

type Frequency string

var Frequencies = [4]string{"WEEKLY", "BIWEEKLY", "TRIWEEKLY", "MONTHLY"}

const (
	WEEKLY    = "WEEKLY"
	BIWEEKLY  = "BIWEEKLY"
	TRIWEEKLY = "TRIWEEKLY"
	MONTHLY   = "MONTHLY"
)
