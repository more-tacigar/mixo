// ==================================================
// Copyright (c) 2016 tacigar. All rights reserved.
// ==================================================

package mixo

import (
	"net/http"
)

// --------------------------------------------------

type methodTrees map[string]*root

type Handler func(*Context)

type Engine struct {
	RouterGroup
	methodTrees methodTrees
}

// --------------------------------------------------

// generate a new engine
func NewEngine() *Engine {
	engine := &Engine{
		RouterGroup: RouterGroup{
			Handlers: []Handler{},
		},
		methodTrees: methodTrees{},
	}
	engine.RouterGroup.engine = engine
	return engine
}

// start mixo application
func (engine *Engine) Run(address string) error {
	err := http.ListenAndServe(address, engine)
	return err
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context := &Context{
		Request:        r,
		ResponseWriter: w,
		engine:         engine,
	}
	root, _ := engine.methodTrees[r.Method]
	relativePath := calculateRelativePath(root.path, r.URL.Path)
	handlers, params := root.getValues(relativePath)

	context.URLParams = params
	for _, handler := range handlers {
		handler(context)
	}
}

// --------------------------------------------------

func (engine *Engine) addRoute(methodName, relativePath string, handlers []Handler) error {
	r, ok := engine.methodTrees[methodName]
	if !ok {
		// if there is not method tree's root, create it.
		r = &root{
			path:     "/",
			handlers: []Handler{},
			children: []*node{},
		}
		engine.methodTrees[methodName] = r
	}
	r.addRoute(relativePath, handlers)
	return nil
}
