package authorization

type Authorization struct {
	ID     int   `json:"id"`
	Status string `json:"status"`
	Data   Data   `json:"data"` 
}

type Data struct {
	Authorization bool `json:"authorization"`
}