[![Go Reference](https://pkg.go.dev/badge/marionette_client.svg)](https://pkg.go.dev/github.com/njasm/marionette_client)
[![CI](https://github.com/njasm/marionette_client/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/njasm/marionette_client/actions/workflows/ci.yml)
[![Coverage Status](https://coveralls.io/repos/github/njasm/marionette_client/badge.svg?branch=master)](https://coveralls.io/github/njasm/marionette_client?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/njasm/marionette_client)](https://goreportcard.com/report/github.com/njasm/marionette_client)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://choosealicense.com/licenses/mit/)

# marionette_client
Mozilla's Gecko Marionette client in golang

## What is Marionette
"Marionette is an automation driver for Mozilla's Gecko engine. It can remotely control either the UI, or the internal 
JavaScript of a Gecko platform, such as Firefox. It can control both the chrome (i.e. menus and functions) or the content 
(the webpage loaded inside the browsing context), giving a high level of control and ability to replicate user actions. 
In addition to performing actions on the browser, Marionette can also read the properties and attributes of the DOM.

If this sounds similar to Selenium/WebDriver then you're correct! Marionette shares much of the same ethos and API as 
Selenium/WebDriver, with additional commands to interact with Gecko's chrome interface. Its goal is to replicate what 
Selenium does for web content: to enable the tester to have the ability to send commands to remotely control a user agent." 

## Resources
https://developer.mozilla.org/en-US/docs/Mozilla/QA/Marionette 

https://w3c.github.io/webdriver/webdriver-spec.html

## Examples
Incomplete list. Check the tests for more examples.

#### Instantiate the client
```go
client := NewClient()
// this are the default marionette values for hostname, and port 
client.Connect("", 0)
// let marionette generate the Session ID with it's default Capabilities
client.NewSession("", nil) 
	
```

#### Navigate to page
```go
client.Navigate("http://www.google.com/")
```

#### Change Contexts
```go
client.SetContext(Context(CHROME))
// or
client.SetContext(Context(CONTENT))
	
```

#### Find Element
```go
element, err := client.FindElement(By(ID), "html-element-id-attribute")
if err != nil {
	// handle your errors
}

// else
println(element.Id())
println(element.Enabled())
println(element.Selected())
println(element.Displayed())
println(element.TagName())
println(element.Text())
println(element.Attribute("id"))
println(element.Property("id"))
println(element.CssValue("text-decoration"))
	
// width, height, x and y
rect, err := element.Rect()
if err != nil {
    // handle your errors
}

fmt.Printf("%#v", rect)
	
// size
w, h, err := element.Size()
if err != nil {
	// handle your errors
}

fmt.Printf("width: %f, height: %f", w, h)

//location
x, y, err := element.Location()
if err != nil {
    // handle your errors
}

fmt.Printf("x: %v, y: %v", x, y)
```

#### Find Elements
```go
collection, err := element.FindElements(By(TAG_NAME), "li")
if err != nil {
	// handle your errors
}

// else
for var e := range collection {
	println(e.Id())
   	println(e.Enabled())
   	println(e.Selected())
   	println(e.Displayed())
   	println(e.TagName())
   	println(e.Text())
   	println(e.Attribute("id"))
   	println(e.CssValue("text-decoration"))
   	e.Click()
}
```

#### Execute JS Script
```go
script := "function mySum(a, b) { return a + b; }; return mySum(arguments[0], arguments[1]);"
args := []int{1, 3} // arguments to be passed to the function
timeout := 1000     // milliseconds
sandbox := false    // new Sandbox
r, err := client.ExecuteScript(script, args, timeout, sandbox)
if err == nil {
    println(r.Value) // 4 
}
```

#### Wait(), Until() Expected condition is true.
```go
client.Navigate("http://www.w3schools.com/ajax/tryit.asp?filename=tryajax_get")

timeout := time.Duration(10) * time.Second
condition := ElementIsPresent(By(ID), "stackH")
ok, webElement, err := Wait(client).For(timeout).Until(condition)

if !ok {
	log.Printf("%#v", err)
	// do your error stuff
	return
}

// cool, we've the element, let's click on it!
webElement.Click()
```
