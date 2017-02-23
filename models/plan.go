package models

import (
	"database/sql"
	"strings"

	"github.com/pborman/uuid"
	"github.com/stripe/stripe-go"
)

/*Plan stores stripe information for a roaster's items*/
type Plan struct {
	RoasterID uuid.UUID      `json:"roasterId"`
	ItemID    uuid.UUID      `json:"itemId"`
	PlanIDs   []string       `json:"planIds"`
	Plans     []*stripe.Plan `json:"plans"`
}

/*PlanRequest contains the information for creating a
  plan in stripe for a roaster */
type PlanRequest struct {
	ItemID uuid.UUID `json:"itemId"`
}

/*NewPlan creates and initializes the ID fields*/
func NewPlan(roasterID, itemID uuid.UUID, planIDs []string) *Plan {
	return &Plan{
		RoasterID: roasterID,
		ItemID:    itemID,
		PlanIDs:   planIDs,
	}
}

/*PlanFromSQL maps sql rows to plan models, where
  order matters*/
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

/*ToFrequency returns the index of the frequency
  string in the plan slices */
func ToFrequency(s Frequency) (int, bool) {
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

/*Frequency is an enum type wrapping string representations
  of the frequency of subscriptions*/
type Frequency string

/*Frequencies are the string representations of Frequency*/
var Frequencies = [4]Frequency{WEEKLY, BIWEEKLY, TRIWEEKLY, MONTHLY}

/*Allowed Frequencies */
const (
	WEEKLY    = "WEEKLY"
	BIWEEKLY  = "BIWEEKLY"
	TRIWEEKLY = "TRIWEEKLY"
	MONTHLY   = "MONTHLY"
)
