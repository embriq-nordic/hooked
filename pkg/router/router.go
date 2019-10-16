package router

import (
	"context"
	"net/http"
	"strings"
)

// Router registers handlers and routes the http request to the right handler.
type Router struct {
	NotFound http.HandlerFunc

	routes map[string]route
}

// New returns a new Router.
func New() *Router {
	return &Router{
		NotFound: http.NotFoundHandler().ServeHTTP,
		routes:   make(map[string]route),
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

// OPTIONS adds the specified handler for OPTIONS requests matching the path.
func (r *Router) OPTIONS(path string, h http.HandlerFunc) {
	r.addHandler(http.MethodOptions, path, h)
}

func segment(path string) []string {
	return strings.Split(strings.Trim(path, "/"), "/")
}

func (r *Router) addHandler(method string, path string, h http.HandlerFunc) {
	rt, exist := r.routes[path]
	if exist {
		rt.handlers[method] = h
		return
	}

	r.routes[path] = route{
		handlers: map[string]http.HandlerFunc{
			method: h,
		},
		segments: segment(path),
	}
}

func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes {
		if ctx, match := route.match(req.Context(), req.URL.Path); match {
			if handler, exist := route.handlers[req.Method]; exist {
				handler.ServeHTTP(res, req.WithContext(ctx))
				return
			}
			var allow []string
			for k := range route.handlers {
				allow = append(allow, k)
			}
			res.Header().Set("Allow", strings.Join(allow, ", "))
			http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
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
	handlers map[string]http.HandlerFunc
	segments []string
}

func (r *route) match(ctx context.Context, path string) (context.Context, bool) {
	pathSegments := segment(path)

	if len(pathSegments) != len(r.segments) {
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
