package marionette_client

type Context int

const (
	Chrome Context = 1 + iota
	Content
)

var contexts = [...]string{
	"chrome",
	"content",
}

// String returns the value name of the context ("chrome", "content").
func (c Context) String() string {
	return contexts[c-1]
}
