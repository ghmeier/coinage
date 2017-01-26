package handlers

import (
	"gopkg.in/alexcesaro/statsd.v2"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/ghmeier/bloodlines/handlers"
	"github.com/jonnykry/coinage/helpers"
)

type CustomerAccountI interface {
	New(ctx *gin.Context)
	View(ctx *gin.Context)
	ViewAll(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type CustomerAccount struct {
	*handlers.BaseHandler
	Helper *helpers.CustomerAccount
}

func NewCustomerAccount(ctx *handlers.GatewayContext) CustomerAccountI {
	stats := ctx.Stats.Clone(statsd.Prefix("api.customer_account"))
	return &CustomerAccount{
		Helper:      helpers.NewCustomerAccount(ctx.Sql),
		BaseHandler: &handlers.BaseHandler{Stats: stats},
	}
}

func (c *CustomerAccount) New(ctx *gin.Context) {
	c.Success(ctx, nil)
}

func (c *CustomerAccount) ViewAll(ctx *gin.Context) {
	c.Success(ctx, nil)
}

func (c *CustomerAccount) View(ctx *gin.Context) {
	c.Success(ctx, nil)
}

func (c *CustomerAccount) Delete(ctx *gin.Context) {
	c.Success(ctx, nil)
}
