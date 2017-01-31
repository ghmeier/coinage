package router

import (
	"fmt"

	"gopkg.in/alexcesaro/statsd.v2"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/ghmeier/bloodlines/config"
	"github.com/ghmeier/bloodlines/gateways"
	h "github.com/ghmeier/bloodlines/handlers"
	"github.com/jonnykry/coinage/handlers"
)

type Billing struct {
	router              *gin.Engine
	roasterAccount      handlers.RoasterAccountI
	customer            handlers.CustomerI
	billingSubscription handlers.BillingSubscriptionI
}

func New(config *config.Root) (*Billing, error) {
	sql, err := gateways.NewSQL(config.SQL)
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

	ctx := &h.GatewayContext{
		Sql:   sql,
		Stats: stats,
	}

	b := &Billing{
		roasterAccount:      handlers.NewRoasterAccount(ctx),
		customer:            handlers.NewCustomer(ctx),
		billingSubscription: handlers.NewBillingSubscription(ctx),
	}
	b.router = gin.Default()
	b.router.Use(h.GetCors())

	roaster := b.router.Group("/api/roaster")
	{
		roaster.POST("", b.roasterAccount.New)
		roaster.GET("", b.roasterAccount.ViewAll)
		roaster.GET("/:accountId", b.roasterAccount.View)
		roaster.PUT("/:accountId", b.roasterAccount.Update)
		roaster.DELETE("/:accountId", b.roasterAccount.Deactivate)
	}
	customer := b.router.Group("/api/customer")
	{
		customer.POST("", b.customer.New)
		customer.GET("", b.customer.ViewAll)
		customer.GET("/:id", b.customer.View)
		customer.DELETE("/:id", b.customer.Delete)
	}
	subscription := b.router.Group("/api/subscription")
	{
		subscription.POST("", b.billingSubscription.New)
		subscription.GET("", b.billingSubscription.Filter)
		subscription.GET("/:subscriptionId", b.billingSubscription.View)
		subscription.PUT("/:subscriptionId", b.billingSubscription.Update)
	}

	return b, nil
}

func (b *Billing) Start(port string) {
	b.router.Run(port)
}
