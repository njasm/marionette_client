package marionette_client

type By int
const (
	ID By = 1 + iota
	CSS_SELECTOR
	LINK_TEXT
	PARTIAL_LINK_TEXT
	XPATH
)
var bys = [...]string {
	"id",
	"css selector",
	"link text",
	"partial link text",
	"xpath",
}

// String returns the English name of the strategy ("css selector", "link text", ...).
func (b By) String() string { return bys[b-1] }

