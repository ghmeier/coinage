package router

import (
	"fmt"

	"gopkg.in/alexcesaro/statsd.v2"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/ghmeier/bloodlines/config"
	g "github.com/ghmeier/bloodlines/gateways"
	h "github.com/ghmeier/bloodlines/handlers"
	w "github.com/ghmeier/bloodlines/workers"
	"github.com/ghmeier/coinage/gateways"
	"github.com/ghmeier/coinage/handlers"
	"github.com/ghmeier/coinage/workers"
	towncenter "github.com/jakelong95/TownCenter/gateways"
	warehouse "github.com/lcollin/warehouse/gateways"
	covenant "github.com/yuderekyu/covenant/gateways"
)

/*Coinage has all the handlers, and routing for the billing microzervice
  README.md has api information*/
type Coinage struct {
	router   *gin.Engine
	roaster  handlers.RoasterI
	customer handlers.CustomerI
	plan     handlers.PlanI
	event    handlers.EventI
	workers  []w.Worker
}

/*New creates and instruments a coinage router*/
func New(config *config.Root) (*Coinage, error) {
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

	rabbit, err := g.NewRabbit(config.Rabbit)
	if err != nil {
		fmt.Println("ERROR: could not connect to RabbitMQ")
		fmt.Println(err.Error())
	}

	stripe := gateways.NewStripe(config.Stripe)
	towncenter := towncenter.NewTownCenter(config.TownCenter)
	warehouse := warehouse.NewWarehouse(config.Warehouse)
	covenant := covenant.NewCovenant(config.Covenant)

	ctx := &h.GatewayContext{
		Sql:        sql,
		Stats:      stats,
		Stripe:     stripe,
		TownCenter: towncenter,
		Warehouse:  warehouse,
		Covenant:   covenant,
		Rabbit:     rabbit,
	}

	b := &Coinage{
		roaster:  handlers.NewRoaster(ctx),
		customer: handlers.NewCustomer(ctx),
		plan:     handlers.NewPlan(ctx),
		event:    handlers.NewEvent(ctx),
		workers:  []w.Worker{workers.NewEvent(ctx)},
	}
	b.router = gin.Default()
	b.router.Use(h.GetCors())

	// id in this case is UserID
	roaster := b.router.Group("/api/roaster")
	{
		roaster.Use(b.roaster.GetJWT())
		roaster.Use(b.roaster.Time())
		roaster.POST("", b.roaster.New)
		roaster.GET("/:id", b.roaster.View)
		//roaster.PUT("/:id", b.roaster.Update)
		roaster.DELETE("/:id", b.roaster.Deactivate)
	}
	plan := b.router.Group("/api/roaster/:id/plan")
	{
		plan.Use(b.plan.GetJWT())
		plan.Use(b.plan.Time())
		plan.GET("", b.plan.ViewAll)
		plan.POST("", b.plan.New)
		plan.GET("/:itemId", b.plan.View)
		//roaster.PUT("/:id/plan/:itemId", b.plan.Update)
		plan.DELETE("/:itemId", b.plan.Delete)
	}
	customer := b.router.Group("/api/customer")
	{
		customer.Use(b.customer.GetJWT())
		customer.Use(b.customer.Time())
		customer.GET("", b.customer.ViewAll)
		customer.GET("/:id", b.customer.View)
		customer.POST("/:id", b.customer.New)
		customer.POST("/:id/subscription", b.customer.Subscribe)
		customer.DELETE("/:id/subscription/:pid", b.customer.Unsubscribe)
		customer.DELETE("/:id", b.customer.Delete)
	}
	event := b.router.Group("/api/event")
	{
		event.Use(b.event.Time())
		event.POST("", b.event.Handle)
	}

	for _, w := range b.workers {
		w.Consume()
	}

	return b, nil
}

/*Start runs the routing engine in gin*/
func (b *Coinage) Start(port string) {
	b.router.Run(port)
}
