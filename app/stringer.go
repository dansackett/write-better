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

	chars := strings.Split(s, "")

	for i := 0; i < len(chars); i++ {
		nodes = append(nodes, &CharNode{
			Index:  i,
			Char:   chars[i],
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

// OpenTag builds a <span> tag to go before the char (project specific)
func OpenTag(procType string, msg string) string {
	return fmt.Sprintf("<span data-placement=\"top\" data-toggle=\"tooltip\" title=\"%s\" class=\"match type-%s\">", msg, procType)
}

// CloseTag builds a <span> tag to go after the char (project specific)
func CloseTag() string {
	return fmt.Sprintf("</span>")
}
