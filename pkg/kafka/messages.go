package kafka

type ProductEvent struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price string `json:"price"`
}
