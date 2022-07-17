package resp

type Base struct {
	Message string      `json:"message"`
	EventID string      `json:"event_id"`
	Data    interface{} `json:"data"`
}
