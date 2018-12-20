package main

import (
	"encoding/json"
)

type (
	// DOMNode represents a DOM node with the smallest set of attributes
	DOMNode struct {
		Tag      string `json:"tag"`
		Text     string `json:"text,omitempty"`
		Attrs    Attrs  `json:"attrs"`
		Children []*DOMNode
	}

	// Attrs just wraps a simple key-value map
	Attrs map[string]string
)

// A sets an attribute to the given DOM node
func (d *DOMNode) A(name, value string) *DOMNode {
	if d.Attrs == nil {
		d.Attrs = Attrs{}
	}
	d.Attrs[name] = value
	return d
}

// Add includes the DOM node as children of this DOMNode
func (d *DOMNode) Add(children ...*DOMNode) *DOMNode {
	d.Children = append(d.Children, children...)
	return d
}

// JSONString return this DOMNode as a JSON node
func (d *DOMNode) JSONString() string {
	buf, err := json.Marshal(d)
	if err != nil {
		panic("Should never happen: " + err.Error())
	}
	return string(buf)
}
