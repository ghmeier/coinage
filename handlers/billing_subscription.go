package handlers

import (
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/ghmeier/bloodlines/gateways"
	"github.com/jonnykry/expresso-billing/helpers"
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
	ctx.JSON(200, empty())
}

func (b *BillingSubscription) Filter(ctx *gin.Context) {
	ctx.JSON(200, empty())
}

func (b *BillingSubscription) View(ctx *gin.Context) {
	ctx.JSON(200, empty())
}

func (b *BillingSubscription) Update(ctx *gin.Context) {
	ctx.JSON(200, empty())
}
