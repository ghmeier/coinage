package router

import (
	"fmt"

	"gopkg.in/gin-gonic/gin.v1"

	"github.com/jonnykry/expresso-billing/handlers"
	"github.com/jonnykry/expresso-billing/gateways"
)

type Billing struct {
	router 	          *gin.Engine
  roasterAccount    handlers.RoasterAccountIfc
}

func New() (*Billing, error) {
	sql, err := gateways.NewSql()

  if err != nil {
		fmt.Println("ERROR: could not connect to mysql.")
		fmt.Println(err.Error())
		return nil, err
	}

	b := &Billing{
		roasterAccount: 	handlers.NewRoasterAccount(sql),
	}
	b.router = gin.Default()

	roasterAccount := b.router.Group("/api/billing/roaster/account")
	{
		roasterAccount.POST("",b.roasterAccount.New)
		roasterAccount.GET("",b.roasterAccount.ViewAll)
		roasterAccount.GET("/:accountId", b.roasterAccount.View)
		roasterAccount.PUT("/:accountId", b.roasterAccount.Update)
		roasterAccount.DELETE("/:accountId", b.roasterAccount.Deactivate)
	}

	return b, nil
}

func (b *Billing) Start(port string) {
	b.router.Run(port)
}
