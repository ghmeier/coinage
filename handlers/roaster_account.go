package handlers

import (
	"github.com/pborman/uuid"
	"gopkg.in/alexcesaro/statsd.v2"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/ghmeier/bloodlines/handlers"
	"github.com/jonnykry/coinage/helpers"
	"github.com/jonnykry/coinage/models"
)

type RoasterAccountI interface {
	New(ctx *gin.Context)
	ViewAll(ctx *gin.Context)
	View(ctx *gin.Context)
	Update(ctx *gin.Context)
	Deactivate(ctx *gin.Context)
}

type RoasterAccount struct {
	*handlers.BaseHandler
	Helper *helpers.RoasterAccount
}

func NewRoasterAccount(ctx *handlers.GatewayContext) RoasterAccountI {
	stats := ctx.Stats.Clone(statsd.Prefix("api.roaster_account"))
	return &RoasterAccount{
		Helper:      helpers.NewRoasterAccount(ctx.Sql),
		BaseHandler: &handlers.BaseHandler{Stats: stats},
	}
}

func (c *RoasterAccount) New(ctx *gin.Context) {
	var json models.RoasterAccount
	err := ctx.BindJSON(&json)

	if err != nil {
		c.UserError(ctx, "Error: unable to parse json", nil)
		return
	}

	// TODO: create stripe account?

	account := models.NewRoasterAccount(json.UserID, json.AccountID)
	err = c.Helper.Insert(account)
	if err != nil {
		c.ServerError(ctx, err, json)
		return
	}

	c.Success(ctx, account)
}

func (c *RoasterAccount) ViewAll(ctx *gin.Context) {
	offset, limit := c.GetPaging(ctx)

	accounts, err := c.Helper.GetAll(offset, limit)
	if err != nil {
		c.ServerError(ctx, err, nil)
		return
	}

	c.Success(ctx, accounts)
}

func (c *RoasterAccount) View(ctx *gin.Context) {
	id := ctx.Param("accountId")

	account, err := c.Helper.GetByID(uuid.Parse(id))
	if err != nil {
		c.ServerError(ctx, err, nil)
		return
	}

	c.Success(ctx, account)
}

func (c *RoasterAccount) Update(ctx *gin.Context) {
	id := ctx.Param("accountId")

	var json models.RoasterAccount
	err := ctx.BindJSON(&json)
	if err != nil {
		c.UserError(ctx, "Error: unable to parse json", nil)
		return
	}

	err = c.Helper.Update(&json)
	if err != nil {
		c.ServerError(ctx, err, json)
		return
	}

	json.ID = uuid.Parse(id)
	c.Success(ctx, json)
}

func (c *RoasterAccount) Deactivate(ctx *gin.Context) {
	c.Success(ctx, nil)
}
