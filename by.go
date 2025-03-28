package marionette_client

import "fmt"

type ByStrategy interface {
	fmt.Stringer
}

// By strategy type to find elements in the DOM
type By int

/*
The method to use to locate the element; one of:

	"id", "name", "class name", "tag name", "css selector",
	"link text", "partial link text", "xpath", "anon" and "anon
	attribute". Note that the "name", "link text" and "partial
	link test" methods are not supported in the chrome DOM.

:param target: The target of the search.  For example, if method =

	"tag", target might equal "div".  If method = "id", target would
	be an element id.

:param id: If specified, search for elements only inside the element

	with the specified id.
*/
const (
	Id By = 1 + iota
	Name
	ClassName
	TagName
	CssSelector
	LinkText
	PartialLinkText
	Xpath
	Anon
	AnonAttribute
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
