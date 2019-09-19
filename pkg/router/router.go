package router

import (
	"context"
	"net/http"
	"strings"
)

// Router registers handlers and routes the http request to the right handler.
type Router struct {
	NotFound http.HandlerFunc

	routes []route
}

// New returns a new Router.
func New() *Router {
	return &Router{
		NotFound: http.NotFoundHandler().ServeHTTP,
	}
}

// GET adds the specified handler for GET requests matching the path.
func (r *Router) GET(path string, h http.HandlerFunc) {
	r.addHandler(http.MethodGet, path, h)
}

// POST adds the specified handler for POST requests matching the path.
func (r *Router) POST(path string, h http.HandlerFunc) {
	r.addHandler(http.MethodPost, path, h)
}

// PUT adds the specified handler for PUT requests matching the path.
func (r *Router) PUT(path string, h http.HandlerFunc) {
	r.addHandler(http.MethodPut, path, h)
}

// DELETE adds the specified handler for DELETE requests matching the path.
func (r *Router) DELETE(path string, h http.HandlerFunc) {
	r.addHandler(http.MethodDelete, path, h)
}

func segment(path string) []string {
	return strings.Split(strings.Trim(path, "/"), "/")
}

func (r *Router) addHandler(method string, path string, h http.HandlerFunc) {
	r.routes = append(r.routes, route{
		method:   method,
		handler:  h,
		segments: segment(path),
	})
}

func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes {
		if ctx, match := route.match(req.Context(), req.Method, req.URL.Path); match {
			route.handler.ServeHTTP(res, req.WithContext(ctx))
			return
		}
	}

	r.NotFound.ServeHTTP(res, req)
}

// paramKey type for adding values to context without risking collision. Eg. two diffenrent packages adding the same string key.
type paramKey string

// GetParam gets a request parameter from context.
func GetParam(ctx context.Context, name string) (string, bool) {
	vStr, ok := ctx.Value(paramKey(name)).(string)
	return vStr, ok
}

type route struct {
	method   string
	handler  http.HandlerFunc
	segments []string
}

func (r *route) match(ctx context.Context, method string, path string) (context.Context, bool) {
	pathSegments := segment(path)

	if (method != r.method) || (len(pathSegments) != len(r.segments)) {
		return nil, false
	}

	for i, seg := range r.segments {
		if strings.HasPrefix(seg, ":") {
			ctx = context.WithValue(ctx, paramKey(strings.Trim(seg, ":")), pathSegments[i])
		} else {
			if seg != pathSegments[i] {
				return nil, false
			}
		}
	}

	return ctx, true
}
