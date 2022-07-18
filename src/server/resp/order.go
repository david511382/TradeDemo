package resp

type OrderPostTest struct {
	Message string `json:"message"`
}

type OrderPostBuy struct {
	Message string `json:"message"`
}

type OrderPostSell struct {
	OrderPostBuy
}
