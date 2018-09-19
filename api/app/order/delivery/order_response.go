package delivery

// Place Order Response
type PlaceOrderResponse struct {
	ID       int         `json:"id"`
	Distance int        `json:"distance"`
	Status   string        `json:"status"`
}

// Take Order Response
type TakeOrderResponse struct {
	Status   string        `json:"status"`
}
