package handlers

import (
	"github.com/pborman/uuid"
	"gopkg.in/alexcesaro/statsd.v2"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/ghmeier/bloodlines/handlers"
	"github.com/jonnykry/coinage/helpers"
	"github.com/jonnykry/coinage/models"
)

type PlanI interface {
	New(ctx *gin.Context)
	View(ctx *gin.Context)
	ViewAll(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type Plan struct {
	*handlers.BaseHandler
	Helper *helpers.Plan
}

func NewPlan(ctx *handlers.GatewayContext) PlanI {
	stats := ctx.Stats.Clone(statsd.Prefix("api.plan"))
	return &Plan{
		Helper:      helpers.NewPlan(ctx.Sql, ctx.Stripe),
		BaseHandler: &handlers.BaseHandler{Stats: stats},
	}
}

func (b *Plan) New(ctx *gin.Context) {
	var json models.PlanRequest
	err := ctx.BindJSON(&json)
	if err != nil {
		b.UserError(ctx, "Error: unable to parse json", err)
		return
	}

	plan, err = b.Helper.Insert(json)
	if err != nil {
		b.ServerError(ctx, err, json)
		return
	}

	b.Success(ctx, plan)
}

func (b *Plan) ViewAll(ctx *gin.Context) {
	id := ctx.Query("id")
	offset, limit := b.GetPaging(ctx)

	plans, err := b.Helper.GetByRoaster(uuid.Parse(id), offset, limit)
	if err != nil {
		b.ServerError(ctx, err, id)
		return
	}

	b.Success(ctx, plans)
}

func (b *Plan) View(ctx *gin.Context) {
	id := ctx.Param("id")
	planID := ctx.Params("pid")

	plan, err := b.Helper.GetPlan(uuid.Parse(id), uuid.Parse(planID))
	if err != nil {
		b.ServerError(ctx, err, nil)
		return
	}

	b.Success(ctx, subscription)
}

func (b *Plan) Update(ctx *gin.Context) {
	var json models.PlanRequests
	err := ctx.BindJSON(&json)
	if err != nil {
		b.UserError(ctx, "Error: unable to parse json", err)
		return
	}

	err = b.Helper.Update(&json)
	if err != nil {
		b.ServerError(ctx, err, json)
		return
	}

	b.Success(ctx, nil)
}
