package handlers

import (
	"fmt"

	//"github.com/pborman/uuid"
	"github.com/stripe/stripe-go"
	"gopkg.in/alexcesaro/statsd.v2"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/ghmeier/bloodlines/handlers"
	"github.com/ghmeier/coinage/helpers"
	"github.com/ghmeier/coinage/workers"
)

/*PlanI describes the requests about billing plans that
  can be handled*/
type EventI interface {
	Handle(*gin.Context)
	/*Time tracks the length of execution for each call in the handler*/
	Time() gin.HandlerFunc
}

/*Plan implements PlanI with coinage helpers*/
type Event struct {
	*handlers.BaseHandler
	Plan    *helpers.Plan
	Roaster *helpers.Roaster
	Event   helpers.Event
}

/*NewEvent initializes and returns a plan with the given gateways*/
func NewEvent(ctx *handlers.GatewayContext) EventI {
	stats := ctx.Stats.Clone(statsd.Prefix("api.event"))
	return &Event{
		Plan:        helpers.NewPlan(ctx.Sql, ctx.Stripe, ctx.Warehouse),
		Roaster:     helpers.NewRoaster(ctx.Sql, ctx.Stripe, ctx.TownCenter),
		Event:       helpers.NewEvent(ctx.Rabbit),
		BaseHandler: &handlers.BaseHandler{Stats: stats},
	}
}

func (e *Event) Handle(ctx *gin.Context) {
	var sEvent stripe.Event
	err := ctx.BindJSON(&sEvent)
	if err != nil {
		fmt.Println(err.Error())
		e.ServerError(ctx, err, nil)
		return
	}

	if !workers.Events[sEvent.Type] {
		e.Success(ctx, nil)
		return
	}

	err = e.Event.Send(&sEvent)
	if err != nil {
		e.ServerError(ctx, err, nil)
		return
	}

	e.Success(ctx, nil)
}

/*
invoice events occur when a new subscription is started:
https://stripe.com/docs/api#invoices
https://stripe.com/docs/api#event_types
invoice.created
invoice.payment_failed
invoice.payment_succeeded
invoice.send
invoice.updated
*/
