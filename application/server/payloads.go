package server

type registerCustomerPayload struct {
	Name  string `json:"name" validate:"required,min=2,max=40"`
	Phone string `json:"phone" validate:"required,min=8,max=11"`
}
