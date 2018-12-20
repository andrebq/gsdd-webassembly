package main

import (
	"encoding/json"
)

type (
	// DOMNode represents a DOM node with the smallest set of attributes
	DOMNode struct {
		Tag      string     `json:"tag"`
		Text     string     `json:"text,omitempty"`
		Attrs    Attrs      `json:"attrs"`
		Children []*DOMNode `json:"children"`
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

// T change the text associated with the given DOMNode
func (d *DOMNode) T(v string) *DOMNode {
	d.Text = v
	return d
}

// H builds a new DOMNode with the given tag
func H(tag string) *DOMNode {
	return &DOMNode{
		Tag: tag,
	}
}

// Div tag
func Div() *DOMNode {
	return H("div")
}

// Text span
func Text(value string) *DOMNode {
	return H("span").T(value)
}
