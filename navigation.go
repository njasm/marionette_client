package marionette_client

type Navigator interface {
	Navigate(url string) (*Response, error)
	PageSource() (*Response, error)
	Title() (string, error)
	Url() (string, error)
	Refresh() error
	Back() error
	Forward() error
}
