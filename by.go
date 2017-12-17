package marionette_client

// strategy type to find elements in the DOM
type By int

/*
   :param method: The method to use to locate the element; one of:
       "id", "name", "class name", "tag name", "css selector",
       "link text", "partial link text", "xpath", "anon" and "anon
       attribute". Note that the "name", "link text" and "partial
       link test" methods are not supported in the chrome DOM.
   :param target: The target of the search.  For example, if method =
       "tag", target might equal "div".  If method = "id", target would
       be an element id.
   :param id: If specified, search for elements only inside the element
       with the specified id.
   """
*/
const (
	ID By = 1 + iota
	NAME
	CLASS_NAME
	TAG_NAME
	CSS_SELECTOR
	LINK_TEXT
	PARTIAL_LINK_TEXT
	XPATH
	ANON
	ANON_ATTRIBUTE
)

var bys = [...]string{
	"id",
	"name",
	"class name",
	"tag name",
	"css selector",
	"link text",
	"partial link text",
	"xpath",
	"anon",
	"anon attribute",
}

// String returns the value name of the strategy ("css selector", "link text", ...).
func (b By) String() string {
	return bys[b-1]
}
