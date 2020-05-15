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

// SaleRow ...
type SaleRow struct {
	DutyID       infra.ObjectID `json:"duty_id"`
	DutyDate     string         `json:"duty_date"`
	DutyQuantity int            `json:"duty_qtd"`

	ID            infra.ObjectID `json:"sale_id"`
	Status        Status         `json:"sale_status"`
	PaymentMethod string         `json:"sale_payment_method"`

	CandyID    infra.ObjectID `json:"candy_id"`
	CandyName  string         `json:"candy_name"`
	CandyPrice int            `json:"candy_price"`

	CustomerID    infra.ObjectID `json:"customer_id"`
	CustomerName  string         `json:"customer_name"`
	CustomerPhone string         `json:"customer_phone"`
}

// DutySale ...
type DutySale struct {
	ID            infra.ObjectID `json:"id"`
	CandyID       infra.ObjectID `json:"candy_id"`
	CandyName     string         `json:"candy_name"`
	CandyPrice    int            `json:"candy_price"`
	CustomerID    infra.ObjectID `json:"customer_id"`
	CustomerName  string         `json:"customer_name"`
	CustomerPhone string         `json:"customer_phone"`
	PaymentMethod string         `json:"payment_method"`
	Status        Status         `json:"status"`
}

// ConsolidatedDuty ...
type ConsolidatedDuty struct {
	ID              infra.ObjectID `json:"id"`
	Date            string         `json:"date"`
	Quantity        int            `json:"quantity"`
	Subtotal        int            `json:"subtotal"`
	PaidAmount      int            `json:"paid_amount"`
	ScheduledAmount int            `json:"scheduled_amount"`
	Sales           []DutySale     `json:"sales"`
}

// ConsolidatedDuties ...
type ConsolidatedDuties map[infra.ObjectID]ConsolidatedDuty
