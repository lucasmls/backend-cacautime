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
