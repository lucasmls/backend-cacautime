package server

type customerPayload struct {
	Name  string `json:"name" validate:"required,min=2,max=40"`
	Phone string `json:"phone" validate:"required,min=8,max=11"`
}

type dutyPayload struct {
	Date          string `json:"date" validate:"required"`
	CandyQuantity int    `json:"candyQuantity" validate:"required,min=1"`
}

type candyPayload struct {
	Name  string `json:"name" validate:"required,min=3,max=40"`
	Price int    `json:"price" validate:"required,min=2"`
}

type registerSalePayload struct {
	CustomerID    int    `json:"customerId" validate:"required,min=1"`
	DutyID        int    `json:"dutyId" validate:"required,min=1"`
	CandyID       int    `json:"candyId" validate:"required,min=1"`
	Status        string `json:"status" validate:"required,oneof=paid not_paid"`
	PaymentMethod string `json:"paymentMethod" validate:"required,oneof=money transfer scheduled"`
}
