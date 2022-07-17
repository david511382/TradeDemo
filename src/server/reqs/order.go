package reqs

type OrderPostBuy struct {
	Quantity uint    `json:"quantity" form:"quantity" binding:"required" uri:"quantity" url:"quantity"`
	Price    float64 `json:"price" form:"price" binding:"required" uri:"price" url:"price"`
}

type OrderPostSell struct {
	OrderPostBuy
}
