package main

import (
	"fmt"
	"strings"
)

// CharNode represents a single node of a string.
type CharNode struct {
	// Index points to the index of the char
	Index int
	// Char is the value being stored
	Char string
	// Before is a list of alterations that need to be inserted before this char
	Before []string
	// After is a list of alterations that need to be inserted after this char
	After []string
}

// CharNodes references a list of CharNode items
type CharNodes []*CharNode

// ToString transforms a list of CharNodes into a string with the before and
// after items mixed in where they must be.
func (nodes CharNodes) ToString() string {
	var str []string

	for _, node := range nodes {
		// Add before items
		if len(node.Before) > 0 {
			for _, before := range node.Before {
				str = append(str, before)
			}
		}

		// Add actual char
		str = append(str, node.Char)

		// Add after items
		if len(node.After) > 0 {
			for _, after := range node.After {
				str = append(str, after)
			}
		}
	}

	return strings.Join(str, "")
}

// ToCharNodes converts a string to a list of CharNode references for processing
func ToCharNodes(s string) CharNodes {
	var nodes CharNodes
	var emptyStrList []string
	chars := []byte(s)

	// @TODO I have learned that when you split a string, there is a mismatch
	// when using strings.Split because of bytes vs string. Because of this,
	// I'm now casting the characters to bytes to get the correct number of
	// characters (so indices match up) but when converting back to a string
	// in building a CharNode, apostrophes (') are not encoded correctly. This
	// needs to be solved.

	for i, char := range chars {
		nodes = append(nodes, &CharNode{
			Index:  i,
			Char:   string(char),
			Before: emptyStrList,
			After:  emptyStrList,
		})
	}

	return nodes
}

// AddBefore adds a string to the before list
func (c *CharNode) AddBefore(s string) {
	c.Before = append(c.Before, s)
}

// AddAfter adds a string to the after list for the char
func (c *CharNode) AddAfter(s string) {
	c.After = append(c.After, s)
}

// SpanTag builds a <span> tag to go before the char (project specific)
func SpanTag(procType string, msg string) string {
	return fmt.Sprintf("<span class=\"type-%s\" data-msg=\"%s\">", procType, msg)
}

// EndSpanTag builds a <span> tag to go after the char (project specific)
func EndSpanTag() string {
	return fmt.Sprintf("</span>")
}
