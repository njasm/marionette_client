package marionette_client

type Navigater interface {
	Get(url string) (*response, error)
}

type navigation struct {
	c   *Client
}
