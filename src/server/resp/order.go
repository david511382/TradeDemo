package resp

type OrderPostBuy struct {
	Message string `json:"message"`
}

type OrderPostSell struct {
	OrderPostBuy
}
