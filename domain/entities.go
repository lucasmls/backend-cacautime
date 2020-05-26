package domain

import (
	"github.com/lucasmls/backend-cacautime/infra"
)

// Customer ...
type Customer struct {
	ID    infra.ObjectID `json:"id"`
	Name  string         `json:"name"`
	Phone string         `json:"phone"`
}

// Duty ...
type Duty struct {
	ID            infra.ObjectID `json:"id"`
	Date          string         `json:"date"`
	CandyQuantity int            `json:"candyQuantity"`
}

// Candy ...
type Candy struct {
	ID    infra.ObjectID `json:"id"`
	Name  string         `json:"name"`
	Price int            `json:"price"`
}

// Sale ...
type Sale struct {
	ID            infra.ObjectID `json:"id"`
	CustomerID    infra.ObjectID `json:"customerId"`
	DutyID        infra.ObjectID `json:"dutyId"`
	CandyID       infra.ObjectID `json:"candyId"`
	Status        Status         `json:"status"`
	PaymentMethod PaymentMethod  `json:"paymentMethod"`
}
