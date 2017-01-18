package handlers

import (
	"github.com/pborman/uuid"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/ghmeier/bloodlines/gateways"
	"github.com/ghmeier/bloodlines/handlers"
	"github.com/jonnykry/coinage/helpers"
	"github.com/jonnykry/coinage/models"
)

type BillingSubscriptionI interface {
	New(ctx *gin.Context)
	Filter(ctx *gin.Context)
	View(ctx *gin.Context)
	Update(ctx *gin.Context)
}

type BillingSubscription struct {
	Helper *helpers.BillingSubscription
}

func NewBillingSubscription(sql gateways.SQL) BillingSubscriptionI {
	return &BillingSubscription{Helper: helpers.NewBillingSubscription(sql)}
}

func (b *BillingSubscription) New(ctx *gin.Context) {
	var json models.BillingSubscription
	err := ctx.BindJSON(&json)
	if err != nil {
		handlers.UserError(ctx, "Error: unable to parse json", err)
		return
	}

	subscription := models.NewBillingSubscription(
		json.UserID,
		json.SubscriptionID,
		json.PlanID,
		json.Amount,
		json.DueAt,
	)

	err = b.Helper.Insert(subscription)
	if err != nil {
		handlers.ServerError(ctx, err, json)
		return
	}

	handlers.Success(ctx, subscription)
}

func (b *BillingSubscription) Filter(ctx *gin.Context) {
	userID := ctx.Query("userId")

	if userID == "" {
		handlers.UserError(ctx, "Error: userId is required", nil)
		return
	}

	subscription, err := b.Helper.GetByUserID(uuid.Parse(userID))
	if err != nil {
		handlers.ServerError(ctx, err, nil)
		return
	}

	handlers.Success(ctx, subscription)
}

func (b *BillingSubscription) View(ctx *gin.Context) {
	id := ctx.Param("subscriptionId")

	subscription, err := b.Helper.GetByID(uuid.Parse(id))
	if err != nil {
		handlers.ServerError(ctx, err, nil)
		return
	}

	handlers.Success(ctx, subscription)
}

func (b *BillingSubscription) Update(ctx *gin.Context) {
	var json models.BillingSubscription
	err := ctx.BindJSON(&json)
	if err != nil {
		handlers.UserError(ctx, "Error: unable to parse json", err)
		return
	}

	err = b.Helper.Update(&json)
	if err != nil {
		handlers.ServerError(ctx, err, json)
		return
	}

	handlers.Success(ctx, nil)
}
