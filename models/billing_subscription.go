package models

import (
	"database/sql"
	"time"

	"github.com/pborman/uuid"
)

type BillingSubscription struct {
	ID             uuid.UUID `json:"id"`
	UserID         uuid.UUID `json:"userId"`
	SubscriptionID string    `json:"subscriptionId"`
	PlanID         string    `json:"planId"`
	Amount         float64   `json:"amount"`
	CreatedAt      time.Time `json:"createdAt"`
	DueAt          time.Time `json:"dueAt"`
}

func NewBillingSubscription(userId uuid.UUID, subscriptionId string, planId string, amount float64, dueAt time.Time) *BillingSubscription {
	return &BillingSubscription{
		ID:             uuid.NewUUID(),
		UserID:         userId,
		SubscriptionID: subscriptionId,
		PlanID:         planId,
		Amount:         amount,
		CreatedAt:      time.Now(),
		DueAt:          dueAt,
	}
}

/*BillingSubscriptionFromSQL returns a subscription splice from sql rows*/
func BillingSubscriptionFromSQL(rows *sql.Rows) ([]*BillingSubscription, error) {
	subscriptions := make([]*BillingSubscription, 0)

	for rows.Next() {
		b := &BillingSubscription{}
		rows.Scan(&b.ID, &b.UserID, &b.SubscriptionID, &b.PlanID, &b.Amount, &b.CreatedAt, b.DueAt)

		subscriptions = append(subscriptions, b)
	}

	return subscriptions, nil
}
