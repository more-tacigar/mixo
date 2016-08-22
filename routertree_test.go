// ==================================================
// Copyright (c) 2016 tacigar. All rights reserved.
// ==================================================

package mixo

import (
	"testing"
)

func TestAddRouteChildrenLength(t *testing.T) {
	r := &root{
		path:     "/",
		children: []*node{},
	}
	h := []Handler{func(c *Context) {}}
	r.addRoute("auth/", h)
	if l := len(r.children); l != 1 {
		t.Errorf("error: r children length must be 1: %d", l)
	}
	r.addRoute("auth/test", h)
	if l := len(r.children); l != 1 {
		t.Errorf("error: r children length must be 1: %d", l)
	}
	if l := len(r.children[0].children); l != 1 {
		t.Errorf("error: r children[0]'s children length must be 1: %d", l)
	}
	r.addRoute("login/", h)
	if l := len(r.children); l != 2 {
		t.Errorf("error: r children length must be 2: %d", l)
	}
}

func TestAddRouteNodeName(t *testing.T) {
	r := &root{
		path:     "/",
		children: []*node{},
	}
	h := []Handler{func(c *Context) {}}
	r.addRoute("auth/", h)
	if c := r.children[0]; c.name != "auth" {
		t.Errorf("error: the name must be \"auth\": %s", c.name)
	}
	r.addRoute("auth/test/", h)
	if c := r.children[0].children[0]; c.name != "test" {
		t.Errorf("error: the name must be \"test\": %s", c.name)
	}
	r.addRoute("auth/aaaa", h)
	if c := r.children[0].children[1]; c.name != "aaaa" {
		t.Errorf("error: the name must be \"aaaa\": %s", c.name)
	}
}

func TestAddRouteNodePath(t *testing.T) {
	r := &root{
		path:     "/",
		children: []*node{},
	}
	h := []Handler{func(c *Context) {}}
	r.addRoute("auth/", h)
	if c := r.children[0]; c.path != "/auth/" {
		t.Errorf("error: the path must be \"/auth/\": %s", c.path)
	}
	r.addRoute("login", h)
	if c := r.children[1]; c.path != "/login" {
		t.Errorf("error: the path must be \"/login\": %s", c.path)
	}
	r.addRoute("auth/test/", h)
	if c := r.children[0].children[0]; c.path != "/auth/test/" {
		t.Errorf("error: the path must be \"/auth/test/\": %s", c.path)
	}
}

func TestGetValuesNormalURL(t *testing.T) {
	r := &root{
		path:     "/",
		children: []*node{},
	}
	h := []Handler{func(c *Context) {}}
	r.addRoute("auth/", h)
	if hs, _ := r.getValues("auth/"); hs == nil {
		t.Errorf("error: must have a handler")
	}
	r.addRoute("auth/test/", h)
	if hs, _ := r.getValues("auth/test/"); hs == nil {
		t.Errorf("error: must have a handler")
	}
	if hs, _ := r.getValues("bug/"); hs != nil {
		t.Errorf("error: must not have a handler")
	}
}

func TestGetValuesRedirect(t *testing.T) {
	r := &root{
		path:     "/",
		children: []*node{},
	}
	h := []Handler{func(c *Context) {}}
	r.addRoute("auth/", h)
	if hs, _ := r.getValues("auth"); hs == nil {
		t.Error("error: must redirect : auth -> auth/")
	}
}

func TestGetValuesWildCard(t *testing.T) {
	r := &root{
		path:     "/",
		children: []*node{},
	}
	h := []Handler{func(c *Context) {}}

	r.addRoute("auth/:name", h)
	if _, ps := r.getValues("auth/john"); ps["name"] != "john" {
		t.Errorf("error: must be john: %s", ps["name"])
	}

	r.addRoute("auth/:name/:action", h)
	if !r.children[0].hasWildChild {
		t.Errorf("error: child must have wild child")
	}
	if !r.children[0].children[0].hasWildChild {
		t.Errorf("error: child must have wild child")
	}
	if _, ps := r.getValues("auth/john/send"); ps["name"] != "john" && ps["action"] != "send" {
		t.Errorf("error: must be john and send: %s, %s", ps["name"], ps["action"])
	}
}
