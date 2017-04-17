package workers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/pborman/uuid"
	"github.com/stripe/stripe-go"

	"github.com/ghmeier/bloodlines/gateways"
	"github.com/ghmeier/bloodlines/handlers"
	bmodels "github.com/ghmeier/bloodlines/models"
	b "github.com/ghmeier/bloodlines/workers"
	cg "github.com/ghmeier/coinage/gateways"
	"github.com/ghmeier/coinage/helpers"
	towncenter "github.com/jakelong95/TownCenter/gateways"
	covenant "github.com/yuderekyu/covenant/gateways"
	cmodels "github.com/yuderekyu/covenant/models"
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
	B        gateways.Bloodlines
	C        covenant.Covenant
	TC       towncenter.TownCenterI
	Stripe   cg.Stripe
	Customer *helpers.Customer
}

func NewEvent(ctx *handlers.GatewayContext) b.Worker {
	worker := &eventWorker{
		RB:       ctx.Rabbit,
		B:        ctx.Bloodlines,
		C:        ctx.Covenant,
		TC:       ctx.TownCenter,
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
	subscribed, err := e.Customer.GetSubscribedFromConnected(customerID)
	if err != nil {
		return err
	}
	if subscribed == nil {
		return fmt.Errorf("No subscribed for connectedId %s", customerID)
	}
	customer, err := e.Customer.GetByCustomerID(subscribed.CustomerID)
	if err != nil {
		return err
	}
	if customer == nil {
		return fmt.Errorf("No customer for customerId %s", subscribed.CustomerID)
	}

	itemID := uuid.Parse(subscription.Plan.Meta["itemId"])
	userID := customer.UserID
	user, err := e.TC.GetUser(userID)
	fmt.Printf("Received invoice for %s, on for %s\n", userID, itemID)
	_, err = e.B.ActivateTrigger("invoice_create", &bmodels.Receipt{
		UserID: userID,
		Values: map[string]string{
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"amount":     strconv.Itoa(int(subscription.Amount)),
			"date":       time.Now().Local().Format("Mon Jan 2 15:04:05 MST 2006"),
		},
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	r := &cmodels.RequestOrder{
		UserID:    userID,
		ItemID:    itemID,
		NextOrder: time.Unix(subscription.Period.End, 0),
	}
	_, err = e.C.NewOrder(r)
	return err
}
