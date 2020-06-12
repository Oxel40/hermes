package modules

// Config stores configuration information
type Config struct {
	Services []Service `json:"services"`
	Clients  []Client  `json:"clients"`
}

// Service ...
type Service struct {
	Name string `json:"name"`
}

// Client ...
type Client struct {
	Name          string   `json:"name"`
	ID            string   `json:"id"`
	Subscriptions []string `json:"subscriptions"`
}
