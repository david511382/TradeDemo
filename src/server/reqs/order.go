package reqs

type OrderPostTest struct {
	RunTimes int `json:"run_times" form:"run_times" binding:"required" uri:"run_times" url:"run_times" default:"50"`
	OrderPostBuy
}

type OrderPostBuy struct {
	Quantity uint    `json:"quantity" form:"quantity" binding:"required" uri:"quantity" url:"quantity" default:"5"`
	Price    float64 `json:"price" form:"price" binding:"required" uri:"price" url:"price" default:"10"`
}

type OrderPostSell struct {
	OrderPostBuy
}
