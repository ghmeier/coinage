package handlers

import (
	"github.com/pborman/uuid"
	"gopkg.in/alexcesaro/statsd.v2"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/ghmeier/bloodlines/handlers"
	"github.com/ghmeier/coinage/helpers"
	"github.com/ghmeier/coinage/models"
)

type RoasterI interface {
	New(ctx *gin.Context)
	View(ctx *gin.Context)
	Update(ctx *gin.Context)
	Deactivate(ctx *gin.Context)
	Time() gin.HandlerFunc
}

type Roaster struct {
	*handlers.BaseHandler
	Roaster *helpers.Roaster
}

func NewRoaster(ctx *handlers.GatewayContext) RoasterI {
	stats := ctx.Stats.Clone(statsd.Prefix("api.roaster_account"))
	return &Roaster{
		Roaster:     helpers.NewRoaster(ctx.Sql, ctx.Stripe, ctx.TownCenter),
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

	account, err := c.Roaster.Insert(&json)
	if err != nil {
		c.ServerError(ctx, err, json)
		return
	}

	c.Success(ctx, account)
}

func (c *Roaster) View(ctx *gin.Context) {
	id := ctx.Param("id")

	account, err := c.Roaster.GetByUserID(uuid.Parse(id))
	if err != nil {
		c.ServerError(ctx, err, nil)
		return
	}

	if account == nil {
		c.UserError(ctx, "ERROR: roaster does not exist", id)
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
