package handlers

import(
	"fmt"

	"gopkg.in/gin-gonic/gin.v1"
	"github.com/pborman/uuid"

	"github.com/jonnykry/expresso-billing/containers"
	"github.com/jonnykry/expresso-billing/gateways"
)

type RoasterAccountIfc interface {
	New(ctx *gin.Context)
	ViewAll(ctx *gin.Context)
	View(ctx *gin.Context)
	Update(ctx *gin.Context)
	Deactivate(ctx *gin.Context)
}

type RoasterAccount struct {
	sql *gateways.Sql
}

func NewRoasterAccount(sql *gateways.Sql) RoasterAccountIfc {
	return &RoasterAccount{sql: sql}
}

func (c *RoasterAccount) New(ctx *gin.Context) {
	var json containers.RoasterAccount
	err := ctx.BindJSON(&json)

	if err != nil {
		ctx.JSON(400, errResponse("Invalid RoasterAccount Object"))
		fmt.Printf("%s",err.Error())
		return
	}

  // TODO:  Update parameters
	err = c.sql.Modify(
		"INSERT INTO roaster_account VALUE(?)",
		uuid.New())
	if err != nil {
		ctx.JSON(500, &gin.H{"error": err, "message": err.Error()})
		return
	}
	ctx.JSON(200, empty())
}

func (c *RoasterAccount) ViewAll(ctx *gin.Context) {
	rows, err := c.sql.Select("SELECT * FROM roaster_account")
	if err != nil {
		 ctx.JSON(500, errResponse(err.Error()))
		 return
	}
	roasterAccount, err := containers.FromSql(rows)
	if err != nil {
		ctx.JSON(500, errResponse(err.Error()))
		return
	}

	ctx.JSON(200, gin.H{"data": roasterAccount})
}

func (c *RoasterAccount) View(ctx *gin.Context) {
	//var json models.RoasterAccount
	id := ctx.Param("accountId")
	if id == "" {
		ctx.JSON(500, errResponse("accountId is a required parameter"))
		return
	}

	rows, err := c.sql.Select("SELECT * FROM roaster_account WHERE id=?", id)
	if err != nil {
		ctx.JSON(500, errResponse(err.Error()))
		return
	}

	roasterAccount, err := containers.FromSql(rows)
	if err != nil {
		ctx.JSON(500, errResponse(err.Error()))
		return
	}

	ctx.JSON(200, gin.H{"data": roasterAccount})
}

func (c *RoasterAccount) Update(ctx *gin.Context) {
	ctx.JSON(200, empty())
}

func (c *RoasterAccount) Deactivate(ctx *gin.Context) {
	ctx.JSON(200, empty())
}
