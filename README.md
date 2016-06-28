[![GoDoc](https://godoc.org/github.com/njasm/marionette_client?status.svg)](https://godoc.org/github.com/njasm/marionette_client)
[![Build Status](https://travis-ci.org/njasm/marionette_client.svg?branch=master)](https://travis-ci.org/njasm/marionette_client)
[![Coverage Status](https://coveralls.io/repos/github/njasm/marionette_client/badge.svg?branch=master)](https://coveralls.io/github/njasm/marionette_client?branch=master)
# marionette_client
Mozilla's Gecko Marionette client in golang

## What is Marionette
"Marionette is an automation driver for Mozilla's Gecko engine. It can remotely control either the UI or the internal 
JavaScript of a Gecko platform, such as Firefox. It can control both the chrome (i.e. menus and functions) or the content 
(the webpage loaded inside the browsing context), giving a high level of control and ability to replicate user actions. 
In addition to performing actions on the browser, Marionette can also read the properties and attributes of the DOM.

If this sounds similar to Selenium/WebDriver then you're correct! Marionette shares much of the same ethos and API as 
Selenium/WebDriver, with additional commands to interact with Gecko's chrome interface. Its goal is to replicate what 
Selenium does for web content: to enable the tester to have the ability to send commands to remotely control a user agent." 

## Resources
https://developer.mozilla.org/en-US/docs/Mozilla/QA/Marionette

## Examples
Incomplete list. Check the tests for more examples.

#### Instantiate the client
```go
	client := NewClient()
	client.Connect("", 0) // this are the default marionette values for hostname, and port 
	client.NewSession("", nil) // let marionette generate the Session ID with it's default Capabilities
	
```

#### Navigate to page
```go
	cliente.Navigate("http://www.google.com/")
```

#### Change Contexts
```go
    client.SetContext(Context(CHROME))
    //or
	client.SetContext(Context(CONTENT))
	
```

#### Execute JS Script
```go
	script := "function mySum(a, b) { return a + b; }; return mySum(arguments[0], arguments[1]);"
	args := []int{1, 3} // arguments to be passed to the function
	timeout := 1000     // milliseconds
	sandbox := false    // new Sandbox
	r, err := client.ExecuteScript(script, args, timeout, sandbox)
	if err == nil {
	    println(r.Value)    // 4 
	}
```

#### Wait(), Until() Expected condition is true.
```go
	client.Navigate("http://www.w3schools.com/ajax/tryit.asp?filename=tryajax_get")
	
	timeout := time.Duration(10) * time.Second
	condition := ElementIsPresent(By(ID), "stackH")
	ok, webElement, _ := Wait(client).For(timeout).Until(condition)

	if !ok {
		log.Printf("%#v", err)
		// do your error stuff
		return
	}

    // cool, we've the element, let's click on it!
	webElement.Click()
	
```