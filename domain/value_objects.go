package domain

import "github.com/lucasmls/backend-cacautime/infra"

// Status ...
type Status string

const (
	// Paid ...
	Paid Status = "paid"
	// NotPaid ...
	NotPaid Status = "not_paid"
)

// PaymentMethod ...
type PaymentMethod string

const (
	// Money ...
	Money PaymentMethod = "money"
	// Transfer ...
	Transfer PaymentMethod = "transfer"
	// Scheduled ...
	Scheduled PaymentMethod = "scheduled"
)

// DutySale ...
type DutySale struct {
	ID            infra.ObjectID `json:"id"`
	Status        Status         `json:"status"`
	PaymentMethod PaymentMethod  `json:"paymentMethod"`

	CandyID    infra.ObjectID `json:"candyId"`
	CandyName  string         `json:"candyName"`
	CandyPrice int            `json:"candyPrice"`

	CustomerID   infra.ObjectID `json:"customerId"`
	CustomerName string         `json:"customerName"`
}

// DutySales ...
type DutySales struct {
	ID       infra.ObjectID `json:"id"`
	Date     string         `json:"date"`
	Quantity int            `json:"quantity"`

	Subtotal        int `json:"subtotal"`
	PaidAmount      int `json:"paidAmount"`
	ScheduledAmount int `json:"scheduledAmount"`

	Sales []DutySale `json:"sales"`
}
