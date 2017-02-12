package handlers

import (
	"github.com/pborman/uuid"
	"gopkg.in/alexcesaro/statsd.v2"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/ghmeier/bloodlines/handlers"
	"github.com/ghmeier/coinage/helpers"
	"github.com/ghmeier/coinage/models"
)

type CustomerI interface {
	New(ctx *gin.Context)
	View(ctx *gin.Context)
	ViewAll(ctx *gin.Context)
	Delete(ctx *gin.Context)
	UpdatePayment(ctx *gin.Context)
	Subscribe(ctx *gin.Context)
	Unsubscribe(ctx *gin.Context)
}

type Customer struct {
	*handlers.BaseHandler
	Customer *helpers.Customer
	Plan     *helpers.Plan
	Roaster  *helpers.Roaster
}

func NewCustomer(ctx *handlers.GatewayContext) CustomerI {
	stats := ctx.Stats.Clone(statsd.Prefix("api.customer"))
	return &Customer{
		Customer:    helpers.NewCustomer(ctx.Sql, ctx.Stripe, ctx.TownCenter, ctx.Covenant),
		Plan:        helpers.NewPlan(ctx.Sql, ctx.Stripe),
		Roaster:     helpers.NewRoaster(ctx.Sql, ctx.Stripe, ctx.TownCenter),
		BaseHandler: &handlers.BaseHandler{Stats: stats},
	}
}

func (c *Customer) New(ctx *gin.Context) {
	var json models.CustomerRequest
	err := ctx.BindJSON(&json)
	if err != nil {
		c.UserError(ctx, "ERROR: unable to parse body", err.Error())
		return
	}

	customer, err := c.Customer.Get(json.UserID)
	if err == nil && customer != nil {
		c.UserError(ctx, "ERROR: customer already exists", nil)
		return
	} else if err != nil {
		c.ServerError(ctx, err, json)
		return
	}

	customer, err = c.Customer.Insert(&json)
	if err != nil {
		c.ServerError(ctx, err, json)
		return
	}

	c.Success(ctx, customer)
}

func (c *Customer) ViewAll(ctx *gin.Context) {
	c.Success(ctx, nil)
}

func (c *Customer) View(ctx *gin.Context) {
	id := ctx.Param("id")

	customer, err := c.Customer.Get(uuid.Parse(id))
	if err != nil {
		c.ServerError(ctx, err, nil)
		return
	}

	c.Success(ctx, customer)
}

func (c *Customer) Subscribe(ctx *gin.Context) {
	id := ctx.Param("id")
	var json models.SubscribeRequest
	err := ctx.BindJSON(&json)
	if err != nil {
		c.UserError(ctx, "ERROR: invalid subscribe request", err)
		return
	}

	roaster, err := c.Roaster.GetByID(json.RoasterID)
	if err != nil {
		c.ServerError(ctx, err, json)
		return
	}

	plan, err := c.Plan.Get(roaster, json.ItemID)
	if err != nil {
		c.ServerError(ctx, err, json)
		return
	}

	err = c.Customer.Subscribe(uuid.Parse(id), plan, json.Frequency)
	if err != nil {
		c.ServerError(ctx, err, json)
		return
	}

	c.Success(ctx, nil)
}

func (c *Customer) Unsubscribe(ctx *gin.Context) {
	c.Success(ctx, nil)
}

func (c *Customer) UpdatePayment(ctx *gin.Context) {
	id := ctx.Param("id")

	var json models.CustomerRequest
	err := ctx.BindJSON(&json)
	if err != nil {
		c.UserError(ctx, "ERROR: token is a required parameter", err)
		return
	}

	err = c.Customer.AddSource(uuid.Parse(id), json.Token)
	if err != nil {
		c.ServerError(ctx, err, &gin.H{"id": id, "token": json.Token})
		return
	}

	c.Success(ctx, nil)
}

func (c *Customer) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.Customer.Delete(uuid.Parse(id))
	if err != nil {
		c.ServerError(ctx, err, nil)
		return
	}

	c.Success(ctx, nil)
}
