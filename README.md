# coinage
A go service for handling Expresso billing

## Quick Start
```bash
$ docker run ghmeier/coinage:latest
```

## Setup
```bash
$ go get github.com/jonnykry/coinage
$ cd $GOPATH/src/github.com/jonnykry/coinage
$ make deps
$ make run
```

## API

### Roaster

### Customer

#### `POST /api/customer` creates a new stripe customer and stores its token.

Example:

*Request:*
```
POST localhost:8081/api/customer
{
	"userId": "<uuid>",
	"token": "<stripe_token>"
}
```

*Response:*
```
{
	"data": {
		"id": <uuid>,
		"userId": <uuid>,
		"customerId": <stripe_customer_id>,
		"subscriptions": {},
		"sources": {},
		"meta": {}
	}
}
```

#### `GET /api/customer?offset=0&limit=20` gets stripe `limit` customers starting at `offset`

Example:

*Request:*
```
GET localhost:8081/api/customer?offset=0&limit=20
```

*Response:*
```
{
	"data": [
		{
			"id": <uuid>,
			"userId": <uuid>,
			"customerId": <stripe_customer_id>,
			"subscriptions": {},
			"sources": {},
			"meta": {}
		}, ...
	]
}
```


#### `GET /api/customer/:id` gets a stripe customer by userId

Example:

*Request:*
```
GET localhost:8081/api/customer/1
```

*Response:*
```
{
	"data": {
		"id": <uuid>,
		"userId": "1",
		"customerId": <stripe_customer_id>,
		"subscriptions": {},
		"sources": {},
		"meta": {}
	}
}
```

#### `DELETE /api/customer/:id` removes the stripe customer

Example:

*Request:*
```
DELETE localhost:8081/api/customer/1
```

*Response:*
```
{
	"data": true
}
```

#### `POST /api/customer/:id/source` updates a user's default payment option by userId

Example:

*Request:*
```
POST localhost:8081/api/customer/1/source
{
	"token": "<stripe_token>"
}
```

*Response:*
```
{
	"data": {
		"id": <uuid>,
		"userId": <uuid>,
		"customerId": <stripe_customer_id>,
		"subscriptions": {},
		"sources": {},
		"meta": {}
	}
}
```

#### `POST /api/customer/:id/subscription` subscribes the user to a new stripe plan

*NOT IMPLEMENTED*

Example:

*Request:*
```
POST localhost:8081/api/customer/1/subscription
```

#### `DELETE /api/customer/:id/subscription/:pid` unsubscribes the user from the plan where id=:pid

*NOT IMPLEMENTED*

Example:

*Request:*
```
POST localhost:8081/api/customer/1/subscription/2
```