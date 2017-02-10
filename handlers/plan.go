package handlers

import (
	"github.com/pborman/uuid"
	"gopkg.in/alexcesaro/statsd.v2"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/ghmeier/bloodlines/handlers"
	"github.com/ghmeier/coinage/helpers"
	"github.com/ghmeier/coinage/models"
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
	Plan    *helpers.Plan
	Roaster *helpers.Roaster
}

func NewPlan(ctx *handlers.GatewayContext) PlanI {
	stats := ctx.Stats.Clone(statsd.Prefix("api.plan"))
	return &Plan{
		Plan:        helpers.NewPlan(ctx.Sql, ctx.Stripe),
		Roaster:     helpers.NewRoaster(ctx.Sql, ctx.Stripe, ctx.TownCenter),
		BaseHandler: &handlers.BaseHandler{Stats: stats},
	}
}

func (p *Plan) New(ctx *gin.Context) {
	id := ctx.Param("id")
	var json models.PlanRequest
	err := ctx.BindJSON(&json)
	if err != nil {
		p.UserError(ctx, "Error: unable to parse json", err)
		return
	}

	roaster, err := p.Roaster.GetByID(uuid.Parse(id))
	if err != nil {
		p.ServerError(ctx, err, json)
		return
	}

	plan, err := p.Plan.Insert(roaster.ID, roaster.AccountID, &json)
	if err != nil {
		p.ServerError(ctx, err, json)
		return
	}

	p.Success(ctx, plan)
}

func (p *Plan) ViewAll(ctx *gin.Context) {
	id := ctx.Query("id")
	offset, limit := p.GetPaging(ctx)

	roaster, err := p.Roaster.GetByUserID(uuid.Parse(id))
	if err != nil {
		p.ServerError(ctx, err, id)
		return
	}

	plans, err := p.Plan.GetByRoaster(roaster, offset, limit)
	if err != nil {
		p.ServerError(ctx, err, id)
		return
	}

	p.Success(ctx, plans)
}

func (p *Plan) View(ctx *gin.Context) {
	id := ctx.Param("id")
	itemID := ctx.Param("itemId")

	roaster, err := p.Roaster.GetByUserID(uuid.Parse(id))
	if err != nil {
		p.ServerError(ctx, err, id)
		return
	}

	plan, err := p.Plan.Get(roaster, uuid.Parse(itemID))
	if err != nil {
		p.ServerError(ctx, err, nil)
		return
	}

	p.Success(ctx, plan)
}

func (p *Plan) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	itemID := ctx.Param("itemId")

	var json models.PlanRequest
	err := ctx.BindJSON(&json)
	if err != nil {
		p.UserError(ctx, "Error: unable to parse json", err)
		return
	}

	//TODO: get item from db

	plan, err := p.Plan.Update(id, uuid.UUID(itemID))
	if err != nil {
		p.ServerError(ctx, err, json)
		return
	}

	p.Success(ctx, plan)
}

func (p *Plan) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	pid := ctx.Param("pid")

	err := p.Plan.Delete(id, pid)
	if err != nil {
		p.ServerError(ctx, err, nil)
		return
	}

	p.Success(ctx, nil)
}
