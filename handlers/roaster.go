package handlers

import (
	"github.com/pborman/uuid"
	"gopkg.in/alexcesaro/statsd.v2"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/ghmeier/bloodlines/handlers"
	"github.com/jonnykry/coinage/helpers"
	"github.com/jonnykry/coinage/models"
)

type RoasterI interface {
	New(ctx *gin.Context)
	ViewAll(ctx *gin.Context)
	View(ctx *gin.Context)
	Update(ctx *gin.Context)
	Deactivate(ctx *gin.Context)
}

type Roaster struct {
	*handlers.BaseHandler
	Helper *helpers.Roaster
}

func NewRoaster(ctx *handlers.GatewayContext) RoasterI {
	stats := ctx.Stats.Clone(statsd.Prefix("api.roaster_account"))
	return &Roaster{
		Helper:      helpers.NewRoaster(ctx.Sql, ctx.Stripe),
		BaseHandler: &handlers.BaseHandler{Stats: stats},
	}
}

func (c *Roaster) New(ctx *gin.Context) {
	var json models.RoasterRequest
	err := ctx.BindJSON(&json)

	if err != nil {
		c.UserError(ctx, "Error: unable to parse json", nil)
		return
	}

	account, err := c.Helper.Insert(&json)
	if err != nil {
		c.ServerError(ctx, err, json)
		return
	}

	c.Success(ctx, account)
}

func (c *Roaster) ViewAll(ctx *gin.Context) {
	offset, limit := c.GetPaging(ctx)

	accounts, err := c.Helper.GetAll(offset, limit)
	if err != nil {
		c.ServerError(ctx, err, nil)
		return
	}

	c.Success(ctx, accounts)
}

func (c *Roaster) View(ctx *gin.Context) {
	id := ctx.Param("id")

	account, err := c.Helper.GetByID(uuid.Parse(id))
	if err != nil {
		c.ServerError(ctx, err, nil)
		return
	}

	c.Success(ctx, account)
}

func (c *Roaster) Update(ctx *gin.Context) {
	c.Success(ctx, nil)
}

func (c *Roaster) Deactivate(ctx *gin.Context) {
	c.Success(ctx, "NOT IMPLEMENTED")
}
