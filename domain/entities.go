package domain

import "github.com/lucasmls/backend-cacautime/infra"

// Customer ...
type Customer struct {
	ID    infra.ObjectID `json:"id"`
	Name  string         `json:"name"`
	Phone string         `json:"phone"`
}
