package router

import (
	"fmt"

	"gopkg.in/gin-gonic/gin.v1"

	"github.com/ghmeier/bloodlines/config"
	"github.com/ghmeier/bloodlines/gateways"
	"github.com/jonnykry/expresso-billing/handlers"
)

type Billing struct {
	router         *gin.Engine
	roasterAccount handlers.RoasterAccountI
}

func New(config *config.Root) (*Billing, error) {
	sql, err := gateways.NewSQL(config.SQL)
	if err != nil {
		fmt.Println("ERROR: could not connect to mysql.")
		fmt.Println(err.Error())
		return nil, err
	}

	b := &Billing{
		roasterAccount: handlers.NewRoasterAccount(sql),
	}
	b.router = gin.Default()

	roaster := b.router.Group("/api/billing/roaster/account")
	{
		roaster.POST("", b.roasterAccount.New)
		roaster.GET("", b.roasterAccount.ViewAll)
		roaster.GET("/:accountId", b.roasterAccount.View)
		roaster.PUT("/:accountId", b.roasterAccount.Update)
		roaster.DELETE("/:accountId", b.roasterAccount.Deactivate)
	}

	return b, nil
}

func (b *Billing) Start(port string) {
	b.router.Run(port)
}
