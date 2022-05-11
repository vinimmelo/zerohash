package transport

type SubscribeRequest struct {
	Type       string   `json:"type"`
	ProductIds []string `json:"product_ids"`
	Channels   []string `json:"channels"`
}

type Subscriber interface {
	Subscribe(request SubscribeRequest)
}
