package marionette_client

type Navigater interface {
	Get(url string) (*response, error)
	GetPageSource() (*response, error)
	GetTitle() (string, error)
	CurrentUrl() (string, error)
	Refresh() error
	Back() error
	Forward() error
}

type navigation struct {
	c *Client
}
