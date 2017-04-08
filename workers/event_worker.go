package workers

import (
	"encoding/json"
	"fmt"

	"github.com/stripe/stripe-go"

	"github.com/ghmeier/bloodlines/gateways"
	"github.com/ghmeier/bloodlines/handlers"
	b "github.com/ghmeier/bloodlines/workers"
	cg "github.com/ghmeier/coinage/gateways"
	"github.com/ghmeier/coinage/helpers"
	covenant "github.com/yuderekyu/covenant/gateways"
)

var Events = map[string]bool{
	"invoice.created":           true,
	"invoice.payment_failed":    true,
	"invoice.payment_succeeded": true,
	"invoice.send":              true,
	"invoice.updated":           true,
}

type eventWorker struct {
	RB       gateways.RabbitI
	C        covenant.Covenant
	Stripe   cg.Stripe
	Customer *helpers.Customer
}

func NewEvent(ctx *handlers.GatewayContext) b.Worker {
	worker := &eventWorker{
		RB:       ctx.Rabbit,
		C:        ctx.Covenant,
		Stripe:   ctx.Stripe,
		Customer: helpers.NewCustomer(ctx.Sql, ctx.Stripe, ctx.TownCenter),
	}

	return &b.BaseWorker{
		HandleFunc: b.HandleFunc(worker.handle),
		RB:         ctx.Rabbit,
	}
}

func (e *eventWorker) handle(body []byte) {
	var event stripe.Event
	err := json.Unmarshal(body, &event)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	switch event.Type {
	case "invoice.created":
		e.invoiceCreate(&event)

	}
}

//invoiceCreate is dispatched when an invoice for a subscription is created
func (e *eventWorker) invoiceCreate(event *stripe.Event) {
	buf, _ := event.Data.Raw.MarshalJSON()

	invoice := &stripe.Invoice{}
	err := invoice.UnmarshalJSON(buf)
	if err != nil {
		fmt.Println("ERROR: unable to parse invoice json")
		fmt.Println(err.Error())
	}

	customerID := event.Data.Obj["customer"].(string)
	for _, item := range invoice.Lines.Values {
		if item.Type != "subscription" {
			continue
		}

		err = e.createOrder(customerID, item)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func (e *eventWorker) createOrder(customerID string, subscription *stripe.InvoiceLine) error {
	customer, err := e.Customer.GetByCustomerID(customerID)
	if err != nil {
		return err
	}
	if customer == nil {
		return fmt.Errorf("No customer for customerId %s", customerID)
	}

	//itemID := subscription.Plan.Meta["itemId"]
	//userID := customer.UserID
	//TODO: send order creation request based on itemID and userID

	return nil
}
