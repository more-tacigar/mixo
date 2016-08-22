// ==================================================
// Copyright (c) 2016 tacigar. All rights reserved.
// ==================================================

package mixo

import (
	"bytes"
)

type node struct {
	path         string
	name         string
	children     []*node
	hasWildChild bool
	handlers     []Handler
}

// root type is alias to node type
type root node

func (n *node) hasChild() bool {
	return len(n.children) != 0
}

func (n *node) hasChildPath(path string) (bool, *node) {
	for _, child := range n.children {
		if child.path == path {
			return true, child
		}
	}
	return false, nil
}

func (n *node) hasChildRelativePath(relativePath string) (bool, *node) {
	fullpath := n.path + relativePath
	hasChild, child := n.hasChildPath(fullpath)
	return hasChild, child
}

func (n *node) hasChildName(name string) (bool, *node) {
	for _, child := range n.children {
		if child.name == name {
			return true, child
		}
	}
	return false, nil
}

func (n *node) addChild(child *node) {
	n.children = append(n.children, child)
}

func isWildCard(runeValue rune) bool {
	return runeValue == ':' || runeValue == '*'

}

// --------------------------------------------------

// add a new route, and register handlers.
func (r *root) addRoute(relativePath string, handlers []Handler) error {
	currentNode := (*node)(r)
	buffer := bytes.NewBufferString("")
	pathBuffer := bytes.NewBufferString(r.path)
	runeRelativePath := []rune(relativePath)

	for i := 0; i < len(runeRelativePath); i++ {
		runeValue := runeRelativePath[i]
		// store buffer as node name
		nodeName := buffer.String()
		buffer.WriteRune(runeValue)

		if runeValue == '/' {
			pathBuffer.WriteString(buffer.String())
			if hasChild, child := currentNode.hasChildName(nodeName); hasChild {
				hasWildChild := i < len(runeRelativePath)-2 && isWildCard(runeRelativePath[i+1])
				currentNode = child
				currentNode.hasWildChild = hasWildChild

			} else {
				hasWildChild := i < len(runeRelativePath)-2 && isWildCard(runeRelativePath[i+1])
				// create new node and connect currentNode
				tmp := buffer.String()
				newNode := &node{
					path:         pathBuffer.String(),
					name:         tmp[:len(tmp)-1],
					children:     []*node{},
					handlers:     []Handler{},
					hasWildChild: hasWildChild,
				}
				currentNode.addChild(newNode)
				currentNode = newNode
			}
			buffer.Reset() // initialize buffer
		}
	}
	// if buffer is not empty. there isn't traling slash.
	if buffer.String() != "" {
		// path buffer currently has parent node path
		pathBuffer.WriteString(buffer.String())
		if hasChild, _ := currentNode.hasChildName(buffer.String()); hasChild {
			// if handler has been registered, return error.
			return ErrAddingRegisteredRoute
		} else {
			newNode := &node{
				path:         pathBuffer.String(),
				name:         buffer.String(),
				children:     []*node{},
				handlers:     handlers,
				hasWildChild: false,
			}
			currentNode.addChild(newNode)
			return nil
		}
	} else {
		// else, current node is target node
		if len(currentNode.handlers) != 0 {
			return ErrAddingRegisteredRoute
		} else {
			currentNode.handlers = handlers
			return nil
		}
	}
}

type URLParams map[string]string

// get handlers and param from a tree
func (r *root) getValues(relativePath string) ([]Handler, URLParams) {
	currentNode := (*node)(r)
	buffer := bytes.NewBufferString("")
	pathBuffer := bytes.NewBufferString(r.path)
	params := URLParams{}
	runeRelativePath := []rune(relativePath)

	for i := 0; i < len(runeRelativePath); i++ {
		runeValue := runeRelativePath[i]
		nodeName := buffer.String()
		buffer.WriteRune(runeValue)

		if runeValue == '/' {
			pathBuffer.WriteString(buffer.String())

			if hasChild, child := currentNode.hasChildName(nodeName); hasChild {
				currentNode = child
			} else {
				if currentNode.hasWildChild {
					// if current node has wild child, child is only one
					key := currentNode.children[0].name[1:] // ignore ':' and '*'.
					tmp := buffer.String()
					params[key] = tmp[:len(tmp)-1]
					currentNode = currentNode.children[0]
				} else {
					return nil, nil
				}
			}
			buffer.Reset()
		}
	}
	// buffer is not empty. there isn't traling slash
	if buffer.String() != "" {
		pathBuffer.WriteString(buffer.String())
		// if has node that doesn't have traling slash, returns it.
		if hasChild, child := currentNode.hasChildName(buffer.String()); hasChild {
			return child.handlers, params
		} else {
			// else redirect or return nil
			if currentNode.hasWildChild {
				// if current node has wild child, child is only one
				key := currentNode.children[0].name[1:] // ignore ':' and '*'.
				tmp := buffer.String()
				params[key] = tmp[:len(tmp)]
				return currentNode.children[0].handlers, params
			} else {
				return nil, nil
			}
		}
	} else {
		return currentNode.handlers, params
	}
}
