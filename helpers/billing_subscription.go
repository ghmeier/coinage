package helpers

import (
	"github.com/pborman/uuid"

	"github.com/ghmeier/bloodlines/gateways"
	"github.com/jonnykry/coinage/models"
)

type BillingSubscription struct {
	*baseHelper
}

func NewBillingSubscription(sql gateways.SQL) *BillingSubscription {
	return &BillingSubscription{
		baseHelper: &baseHelper{sql: sql},
	}
}

func (b *BillingSubscription) Insert(subscription *models.BillingSubscription) error {
	err := b.sql.Modify("INSERT INTO billing_subscription (id,userId,subscriptionId,planId,amount,createdAt,dueAt)VALUES(?,?,?,?,?,?,?)",
		subscription.ID,
		subscription.UserID,
		subscription.SubscriptionID,
		subscription.PlanID,
		subscription.Amount,
		subscription.CreatedAt,
		subscription.DueAt,
	)
	return err
}

func (b *BillingSubscription) GetByID(id uuid.UUID) (*models.BillingSubscription, error) {
	return b.getOne(
		"SELECT id,userId,subscriptionId,planId,amount,createdAt,dueAt FROM billing_subscription WHERE id=?",
		id,
	)
}

func (b *BillingSubscription) GetByUserID(id uuid.UUID) (*models.BillingSubscription, error) {
	return b.getOne(
		"SELECT id,userId,subscriptionId,planId,amount,createdAt,dueAt FROM billing_subscription WHERE userId=?",
		id,
	)
}

func (b *BillingSubscription) Update(subscription *models.BillingSubscription) error {
	err := b.sql.Modify("UPDATE billing_subscription SET userId=?,subscriptionId=?,planId=?,amount=?,dueAt=? WHERE id=?",
		subscription.UserID,
		subscription.SubscriptionID,
		subscription.PlanID,
		subscription.Amount,
		subscription.DueAt,
	)
	return err
}

func (b *BillingSubscription) getOne(query string, id uuid.UUID) (*models.BillingSubscription, error) {
	rows, err := b.sql.Select(query, id)

	if err != nil {
		return nil, err
	}

	subscriptions, err := models.BillingSubscriptionFromSQL(rows)
	if err != nil {
		return nil, err
	}

	return subscriptions[0], nil
}
