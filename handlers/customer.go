package handlers

import (
	"github.com/pborman/uuid"
	"gopkg.in/alexcesaro/statsd.v2"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/ghmeier/bloodlines/handlers"
	"github.com/ghmeier/coinage/helpers"
	"github.com/ghmeier/coinage/models"
)

/*CustomerI describes the customer handler interface*/
type CustomerI interface {
	/*New creates a new customer, sending it back on success*/
	New(ctx *gin.Context)
	/*View sends the customer with the given ID as a response*/
	View(ctx *gin.Context)
	/*View all sends a list of customers*/
	ViewAll(ctx *gin.Context)
	/*Delete removes the customer with the given ID*/
	Delete(ctx *gin.Context)
	/*UpdatePayment creates and sets a new default payment for the customer*/
	UpdatePayment(ctx *gin.Context)
	/*Subscribe creates a new subscription for the cusutomer*/
	Subscribe(ctx *gin.Context)
	/*Unsubscribe removes a sucbsription from a customer*/
	Unsubscribe(ctx *gin.Context)
	/*Time tracks the length of execution for each call in the handler*/
	Time() gin.HandlerFunc
}

/*Customer implements CustomerI using stripe and coinage helpers*/
type Customer struct {
	*handlers.BaseHandler
	Customer *helpers.Customer
	Plan     *helpers.Plan
	Roaster  *helpers.Roaster
}

/*NewCustomer creates and returns a new customer using the given context*/
func NewCustomer(ctx *handlers.GatewayContext) CustomerI {
	stats := ctx.Stats.Clone(statsd.Prefix("api.customer"))
	return &Customer{
		Customer:    helpers.NewCustomer(ctx.Sql, ctx.Stripe, ctx.TownCenter, ctx.Covenant),
		Plan:        helpers.NewPlan(ctx.Sql, ctx.Stripe),
		Roaster:     helpers.NewRoaster(ctx.Sql, ctx.Stripe, ctx.TownCenter),
		BaseHandler: &handlers.BaseHandler{Stats: stats},
	}
}

/*New implements CustomerI.New*/
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

/*ViewAll implements CustomerI.ViewAll*/
func (c *Customer) ViewAll(ctx *gin.Context) {
	c.Success(ctx, nil)
}

/*View implments CustomerI.View*/
func (c *Customer) View(ctx *gin.Context) {
	id := ctx.Param("id")

	customer, err := c.Customer.Get(uuid.Parse(id))
	if err != nil {
		c.ServerError(ctx, err, nil)
		return
	}

	if customer == nil {
		c.UserError(ctx, "ERROR: customer does not exist", nil)
		return
	}

	c.Success(ctx, customer)
}

/*Subscribe implements CustomerI.Subscribe*/
func (c *Customer) Subscribe(ctx *gin.Context) {
	id := ctx.Param("id")
	var json models.SubscribeRequest
	err := ctx.BindJSON(&json)
	if err != nil {
		c.UserError(ctx, "ERROR: invalid subscribe request", err)
		return
	}

	roaster, err := c.Roaster.Get(json.RoasterID)
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

/*Unsubscribe is not implemented*/
func (c *Customer) Unsubscribe(ctx *gin.Context) {
	c.Success(ctx, nil)
}

/*UpdatePayment implements CustomerI.UpdatePayment*/
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

/*Delete implements CustomerI.Delete*/
func (c *Customer) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.Customer.Delete(uuid.Parse(id))
	if err != nil {
		c.ServerError(ctx, err, nil)
		return
	}

	c.Success(ctx, nil)
}
