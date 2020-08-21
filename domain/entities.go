package domain

import (
	"github.com/lucasmls/backend-cacautime/infra"
)

// User ...
type User struct {
	ID       infra.ObjectID `json:"id"`
	Name     string         `json:"name"`
	Email    string         `json:"email"`
	Password string         `json:"password"`
}

// Customer ...
type Customer struct {
	ID    infra.ObjectID `json:"id"`
	Name  string         `json:"name"`
	Phone string         `json:"phone"`
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
	CandyID       infra.ObjectID `json:"candyId"`
	Status        Status         `json:"status"`
	PaymentMethod PaymentMethod  `json:"paymentMethod"`
	Date          string         `json:"date"`
}

// Month ...
type Month struct {
	Month  string `json:"month"`
	Year   string `json:"year"`
	Number string `json:"number"`
}

// MonthSale ...
type MonthSale struct {
	ID            infra.ObjectID `json:"id"`
	Status        Status         `json:"status"`
	PaymentMethod PaymentMethod  `json:"paymentMethod"`

	CandyID    infra.ObjectID `json:"candyId"`
	CandyName  string         `json:"candyName"`
	CandyPrice int            `json:"candyPrice"`

	CustomerID   infra.ObjectID `json:"customerId"`
	CustomerName string         `json:"customerName"`
}

// MonthSales ...
type MonthSales struct {
	Subtotal        int `json:"subtotal"`
	PaidAmount      int `json:"paidAmount"`
	ScheduledAmount int `json:"scheduledAmount"`

	Sales []MonthSale `json:"sales"`
}