package marionette_client

type Context int
const (
	CHROME Context = 1 + iota
	CONTENT
)

var contexts = [...]string {
	"chrome",
	"content",
}

// String returns the English name of the month ("chrome", "content").
func (c Context) String() string { return contexts[c-1] }
