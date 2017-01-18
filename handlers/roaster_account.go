package handlers

import (
	"github.com/pborman/uuid"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/ghmeier/bloodlines/gateways"
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
	Helper *helpers.RoasterAccount
}

func NewRoasterAccount(sql gateways.SQL) RoasterAccountI {
	return &RoasterAccount{Helper: helpers.NewRoasterAccount(sql)}
}

func (c *RoasterAccount) New(ctx *gin.Context) {
	var json models.RoasterAccount
	err := ctx.BindJSON(&json)

	if err != nil {
		handlers.UserError(ctx, "Error: unable to parse json", nil)
		return
	}

	// TODO: create stripe account?

	account := models.NewRoasterAccount(json.UserID, json.AccountID)
	err = c.Helper.Insert(account)
	if err != nil {
		handlers.ServerError(ctx, err, json)
		return
	}

	handlers.Success(ctx, account)
}

func (c *RoasterAccount) ViewAll(ctx *gin.Context) {
	offset, limit := handlers.GetPaging(ctx)

	accounts, err := c.Helper.GetAll(offset, limit)
	if err != nil {
		handlers.ServerError(ctx, err, nil)
		return
	}

	handlers.Success(ctx, accounts)
}

func (c *RoasterAccount) View(ctx *gin.Context) {
	id := ctx.Param("accountId")

	account, err := c.Helper.GetByID(uuid.Parse(id))
	if err != nil {
		handlers.ServerError(ctx, err, nil)
		return
	}

	handlers.Success(ctx, account)
}

func (c *RoasterAccount) Update(ctx *gin.Context) {
	id := ctx.Param("accountId")

	var json models.RoasterAccount
	err := ctx.BindJSON(&json)
	if err != nil {
		handlers.UserError(ctx, "Error: unable to parse json", nil)
		return
	}

	err = c.Helper.Update(&json)
	if err != nil {
		handlers.ServerError(ctx, err, json)
		return
	}

	json.ID = uuid.Parse(id)
	handlers.Success(ctx, json)
}

func (c *RoasterAccount) Deactivate(ctx *gin.Context) {
	handlers.Success(ctx, nil)
}
