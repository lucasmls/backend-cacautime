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
	PaymentMethod PaymentMethod  `json:"payment_method"`

	CandyID    infra.ObjectID `json:"candy_id"`
	CandyName  string         `json:"candy_name"`
	CandyPrice int            `json:"candy_price"`

	CustomerID   infra.ObjectID `json:"customer_id"`
	CustomerName string         `json:"customer_name"`
}

// DutySales ...
type DutySales struct {
	ID       infra.ObjectID `json:"id"`
	Date     string         `json:"date"`
	Quantity int            `json:"quantity"`

	Subtotal        int `json:"subtotal"`
	PaidAmount      int `json:"paid_amount"`
	ScheduledAmount int `json:"scheduled_amount"`

	Sales []DutySale `json:"sales"`
}
