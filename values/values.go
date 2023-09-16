package values

type LogLine struct {
	ID   int    `json:"id"`
	Timestamp string `json:"timestamp"`
}


type Response struct {
	Counter int `json:"counter"`
}