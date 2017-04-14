package helpers

import (
	"github.com/stripe/stripe-go"

	g "github.com/ghmeier/bloodlines/gateways"
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
