package network

// InMsg received from RabbitMQ
type InMsg struct {
	Exchange   string
	RoutingKey string
	Headers    map[string]interface{}
	Body       []byte
}
