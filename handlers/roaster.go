package handlers

import (
	"github.com/pborman/uuid"
	"gopkg.in/alexcesaro/statsd.v2"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/ghmeier/bloodlines/handlers"
	"github.com/ghmeier/coinage/helpers"
	"github.com/ghmeier/coinage/models"
)

/*RoasterI describes the roaster requests that Coinage can handle*/
type RoasterI interface {
	/*New creates a new roaster billing account, sending the result*/
	New(ctx *gin.Context)
	/*View sends the roaster account with the given ID*/
	View(ctx *gin.Context)
	/*Update sets the propperties of a roaster account, sending the update*/
	Update(ctx *gin.Context)
	/*Deactivate removes a roaster account from coinage*/
	Deactivate(ctx *gin.Context)
	/*Time tracks the duration of each handled request*/
	Time() gin.HandlerFunc
	/*GetJWT checks the JWT of a request*/
	GetJWT() gin.HandlerFunc
}

/*Roaster implements RoasterI*/
type Roaster struct {
	*handlers.BaseHandler
	Roaster *helpers.Roaster
}

/*NewRoaster returns a new roaster with the given gateways*/
func NewRoaster(ctx *handlers.GatewayContext) RoasterI {
	stats := ctx.Stats.Clone(statsd.Prefix("api.roaster"))
	return &Roaster{
		Roaster:     helpers.NewRoaster(ctx.Sql, ctx.Stripe, ctx.TownCenter),
		BaseHandler: &handlers.BaseHandler{Stats: stats},
	}
}

/*New implements RoasterI.New*/
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

/*View implements RoasterI.View*/
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

/*Update is not implemented*/
func (c *Roaster) Update(ctx *gin.Context) {
	c.Success(ctx, "NOT IMPLMENTED")
}

/*Deactivate is not implemented*/
func (c *Roaster) Deactivate(ctx *gin.Context) {
	c.Success(ctx, "NOT IMPLEMENTED")
}
