package helpers

import (
	//	"fmt"

	//	"github.com/pborman/uuid"
	"github.com/stripe/stripe-go"

	g "github.com/ghmeier/bloodlines/gateways"
	//"github.com/ghmeier/coinage/gateways"
	//	"github.com/ghmeier/coinage/models"
	//towncenter "github.com/jakelong95/TownCenter/gateways"
	//	t "github.com/jakelong95/TownCenter/models"
	//warehouse "github.com/lcollin/warehouse/gateways"
	//covenant "github.com/yuderekyu/covenant/gateways"
)

/*Event helps with manipulating roaster properties*/
type Event interface {
	Send(*stripe.Event) error
}

type event struct {
	R g.RabbitI
}

/*NewEvent initializes and returns a roaster with the given gateways*/
func NewEvent(rabbit g.RabbitI) Event {
	return &event{
		R: rabbit,
	}
}

func (e *event) Send(sEvent *stripe.Event) error {
	return e.R.Produce(sEvent)
}
