// ==================================================
// Copyright (c) 2016 tacigar. All rights reserved.
// ==================================================

package mixo

type RouterGroup struct {
	engine   *Engine
	Handlers []Handler
}

// --------------------------------------------------

func (routerGroup *RouterGroup) registerHandler(
	methodName, relativePath string, handlers []Handler) {
	// append router group's handlers and parameter handlers
	newHandlers := append(routerGroup.Handlers, handlers...)
	routerGroup.engine.addRoute(methodName, relativePath, newHandlers)
}

func (routerGroup *RouterGroup) GET(relativePath string, handlers ...Handler) {
	routerGroup.registerHandler("GET", relativePath, handlers)
}

func (routerGroup *RouterGroup) POST(relativePath string, handlers ...Handler) {
	routerGroup.registerHandler("POST", relativePath, handlers)
}

func (routerGroup *RouterGroup) PUT(relativePath string, handlers ...Handler) {
	routerGroup.registerHandler("PUT", relativePath, handlers)
}

func (routerGroup *RouterGroup) DELETE(relativePath string, handlers ...Handler) {
	routerGroup.registerHandler("DELETE", relativePath, handlers)
}

func (routerGroup *RouterGroup) Use(handlers ...Handler) {
	routerGroup.Handlers = append(routerGroup.Handlers, handlers...)
}

// --------------------------------------------------

func (routerGroup *RouterGroup) Branch() *RouterGroup {
	return &RouterGroup{
		engine:   routerGroup.engine,
		Handlers: routerGroup.Handlers,
	}
}
