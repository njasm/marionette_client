package marionette_client

type Navigator interface {
	Navigate(url string) (*response, error)
	PageSource() (*response, error)
	Title() (string, error)
	CurrentUrl() (string, error)
	Refresh() error
	Back() error
	Forward() error
}
