package marionette_client

type Navigator interface {
	Navigate(url string) (*response, error)
	GetPageSource() (*response, error)
	GetTitle() (string, error)
	CurrentUrl() (string, error)
	Refresh() error
	Back() error
	Forward() error
}
