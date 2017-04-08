package models

import (
	"database/sql"
	"strings"

	"github.com/pborman/uuid"
)

/*Plan stores stripe information for a roaster's items*/
type Plan struct {
	RoasterID uuid.UUID `json:"roasterId"`
	ItemID    uuid.UUID `json:"itemId"`
	PlanIDs   []string  `json:"planIds"`
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
		return 1, true
	case BIWEEKLY:
		return 2, true
	case TRIWEEKLY:
		return 3, true
	case MONTHLY:
		return 4, true
	default:
		return 0, false
	}
}

/*Frequency is an enum type wrapping string representations
  of the frequency of subscriptions*/
type Frequency string

/*Frequencies are the string representations of Frequency*/
var Frequencies = [5]Frequency{INVALID, WEEKLY, BIWEEKLY, TRIWEEKLY, MONTHLY}

/*Allowed Frequencies */
const (
	INVALID   = "INVALID"
	WEEKLY    = "WEEKLY"
	BIWEEKLY  = "BIWEEKLY"
	TRIWEEKLY = "TRIWEEKLY"
	MONTHLY   = "MONTHLY"
)
