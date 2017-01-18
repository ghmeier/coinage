package handlers

import (
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/ghmeier/bloodlines/gateways"
	"github.com/jonnykry/expresso-billing/helpers"
)

type CustomerAccountI interface {
	New(ctx *gin.Context)
	View(ctx *gin.Context)
	ViewAll(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type CustomerAccount struct {
	Helper *helpers.CustomerAccount
}

func NewCustomerAccount(sql gateways.SQL) CustomerAccountI {
	return &CustomerAccount{
		Helper: helpers.NewCustomerAccount(sql),
	}
}

func (c *CustomerAccount) New(ctx *gin.Context) {
	ctx.JSON(200, empty())
}

func (c *CustomerAccount) ViewAll(ctx *gin.Context) {
	ctx.JSON(200, empty())
}

func (c *CustomerAccount) View(ctx *gin.Context) {
	ctx.JSON(200, empty())
}

func (c *CustomerAccount) Delete(ctx *gin.Context) {
	ctx.JSON(200, empty())
}
