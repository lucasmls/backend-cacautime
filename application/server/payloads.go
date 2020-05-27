package server

type registerCustomerPayload struct {
	Name  string `json:"name" validate:"required,min=2,max=40"`
	Phone string `json:"phone" validate:"required,min=8,max=11"`
}

type registerDutyPayload struct {
	Date          string `json:"date" validate:"required"`
	CandyQuantity int    `json:"candyQuantity" validate:"required,min=1"`
}

type registerCandyPayload struct {
	Name  string `json:"name" validate:"required,min=3,max=40"`
	Price int    `json:"price" validate:"required,min=2"`
}
