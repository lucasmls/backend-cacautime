package server

type loginPayload struct {
	Email    string `json:"email" validate:"required,min=2,max=40"`
	Password string `json:"password" validate:"required,min=3,max=100"`
}

type customerPayload struct {
	Name  string `json:"name" validate:"required,min=2,max=40"`
	Phone string `json:"phone" validate:"max=11"`
}

type candyPayload struct {
	Name  string `json:"name" validate:"required,min=3,max=100"`
	Price int    `json:"price" validate:"required,min=2"`
}

type salePayload struct {
	CustomerID    int    `json:"customerId" validate:"required,min=1"`
	CandyID       int    `json:"candyId" validate:"required,min=1"`
	Status        string `json:"status" validate:"required,oneof=paid not_paid"`
	PaymentMethod string `json:"paymentMethod" validate:"required,oneof=money transfer scheduled"`
}

type updateSalePayload struct {
	Status        string `json:"status" validate:"required,oneof=paid not_paid"`
	PaymentMethod string `json:"paymentMethod" validate:"required,oneof=money transfer scheduled"`
}
