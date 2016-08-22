// ==================================================
// Copyright (c) 2016 tacigar. All rights reserved.
// ==================================================

package mixo

import (
	"html/template"
	"net/http"
)

type Context struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	URLParams      URLParams
	engine         *Engine
	Metadata       map[string]interface{}
}

// --------------------------------------------------

func (context *Context) SetMetadata(key string, value interface{}) {
	if context.Metadata == nil {
		context.Metadata = make(map[string]interface{})
	}
	context.Metadata[key] = value
}

func (context *Context) GetMetadata(key string) (interface{}, bool) {
	if context.Metadata == nil {
		return nil, false
	}
	value, ok := context.Metadata[key]
	return value, ok
}

// --------------------------------------------------

func (context *Context) Param(key string) string {
	return context.URLParams[key]
}

// get query parameter by key
func (context *Context) Query(key string) (string, bool) {
	request := context.Request
	if values, ok := request.URL.Query()[key]; ok && len(values) > 0 {
		return values[0], true
	}
	return "", false
}

// return post form value
func (context *Context) PostForm(key string) string {
	req := context.Request
	req.ParseForm()
	if values := req.PostForm[key]; len(values) > 0 {
		return values[0]
	}
	return ""
}

func (context *Context) RenderHTML(
	code int, fileName string, data map[string]interface{}) error {

	temp, err := template.ParseFiles(fileName)
	if err != nil {
		return err
	}
	context.ResponseWriter.Header().Set("Content-Type", "text/html; charset-utf-8")
	context.ResponseWriter.WriteHeader(code)
	err = temp.Execute(context.ResponseWriter, data)
	if err != nil {
		return err
	}
	return nil
}

func (context *Context) AbortWithStatus(code int) {
	context.ResponseWriter.WriteHeader(code)
}
