//+build js
package main

import (
	"syscall/js"
)

// Render the given node to the root object
func Render(n *DOMNode) {
	js.Global().Call("renderJSON", n.JSONString())
}
