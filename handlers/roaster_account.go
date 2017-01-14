package handlers

import (
	"fmt"

	"github.com/pborman/uuid"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/ghmeier/bloodlines/gateways"
	"github.com/jonnykry/expresso-billing/helpers"
	"github.com/jonnykry/expresso-billing/models"
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
		ctx.JSON(400, errResponse("Invalid RoasterAccount Object"))
		fmt.Printf("%s", err.Error())
		return
	}

	// TODO: create stripe account?

	account := models.NewRoasterAccount(json.UserID, json.AccountID)
	err = c.Helper.Insert(account)
	if err != nil {
		ctx.JSON(500, &gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, &gin.H{"data": account})
}

func (c *RoasterAccount) ViewAll(ctx *gin.Context) {
	offset, limit := getPaging(ctx)

	accounts, err := c.Helper.GetAll(offset, limit)
	if err != nil {
		ctx.JSON(500, errResponse(err.Error()))
		return
	}

	ctx.JSON(200, gin.H{"data": accounts})
}

func (c *RoasterAccount) View(ctx *gin.Context) {
	id := ctx.Param("accountId")

	account, err := c.Helper.GetByID(uuid.Parse(id))
	if err != nil {
		ctx.JSON(500, errResponse(err.Error()))
		return
	}

	ctx.JSON(200, gin.H{"data": account})
}

func (c *RoasterAccount) Update(ctx *gin.Context) {
	id := ctx.Param("accountId")

	var json models.RoasterAccount
	err := ctx.BindJSON(&json)
	if err != nil {
		ctx.JSON(400, errResponse(err.Error()))
		return
	}

	err = c.Helper.Update(&json)
	if err != nil {
		ctx.JSON(500, errResponse(err.Error()))
		return
	}

	json.ID = uuid.Parse(id)
	ctx.JSON(200, &gin.H{"data": json})
}

func (c *RoasterAccount) Deactivate(ctx *gin.Context) {
	ctx.JSON(200, empty())
}
