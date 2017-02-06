package router

import (
	"fmt"

	"gopkg.in/alexcesaro/statsd.v2"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/ghmeier/bloodlines/config"
	g "github.com/ghmeier/bloodlines/gateways"
	h "github.com/ghmeier/bloodlines/handlers"
	towncenter "github.com/jakelong95/TownCenter/gateways"
	"github.com/jonnykry/coinage/gateways"
	"github.com/jonnykry/coinage/handlers"
)

type Billing struct {
	router   *gin.Engine
	roaster  handlers.RoasterI
	customer handlers.CustomerI
	plan     handlers.PlanI
}

func New(config *config.Root) (*Billing, error) {
	sql, err := g.NewSQL(config.SQL)
	if err != nil {
		fmt.Println("ERROR: could not connect to mysql.")
		fmt.Println(err.Error())
	}

	stats, err := statsd.New(
		statsd.Address(config.Statsd.Host+":"+config.Statsd.Port),
		statsd.Prefix(config.Statsd.Prefix),
	)
	if err != nil {
		fmt.Println("ERROR: unable to connect to statsd")
		fmt.Println(err.Error())
	}

	stripe := gateways.NewStripe(config.Stripe)
	towncenter := towncenter.NewTownCenter(config.TownCenter)

	ctx := &h.GatewayContext{
		Sql:        sql,
		Stats:      stats,
		Stripe:     stripe,
		TownCenter: towncenter,
	}

	b := &Billing{
		roaster:  handlers.NewRoaster(ctx),
		customer: handlers.NewCustomer(ctx),
		plan:     handlers.NewPlan(ctx),
	}
	b.router = gin.Default()
	b.router.Use(h.GetCors())

	roaster := b.router.Group("/api/roaster")
	{
		roaster.POST("", b.roaster.New)
		roaster.GET("", b.roaster.ViewAll)
		roaster.GET("/:id", b.roaster.View)
		//roaster.PUT("/:id", b.roaster.Update)
		roaster.DELETE("/:id", b.roaster.Deactivate)
		roaster.GET("/:id/plan", b.plan.ViewAll)
		roaster.POST("/:id/plan", b.plan.New)
		roaster.GET("/:id/plan/:pid", b.plan.View)
		roaster.PUT("/:id/plan/:pid", b.plan.Update)
		roaster.DELETE("/:id/plan/:pid", b.plan.Delete)
	}
	customer := b.router.Group("/api/customer")
	{
		customer.POST("", b.customer.New)
		customer.GET("", b.customer.ViewAll)
		customer.GET("/:id", b.customer.View)
		customer.POST("/:id/source", b.customer.UpdatePayment)
		customer.POST("/:id/subscription", b.customer.Subscribe)
		customer.DELETE("/:id", b.customer.Delete)
	}

	return b, nil
}

func (b *Billing) Start(port string) {
	b.router.Run(port)
}
