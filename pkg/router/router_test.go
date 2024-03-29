package router

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var serveTests = []struct {
	method          string
	path            string
	expectedBody    string
	expectedHeaders map[string]string
}{
	{
		method:       http.MethodGet,
		path:         "/invalidpath",
		expectedBody: "NotFound",
	},
	{
		method:       http.MethodGet,
		path:         "/test",
		expectedBody: "GETTest",
	},
	{
		method:       http.MethodPost,
		path:         "/test",
		expectedBody: "POSTTest",
	},
	{
		method:       http.MethodPut,
		path:         "/test",
		expectedBody: "PUTTest",
	},
	{
		method:       http.MethodDelete,
		path:         "/test",
		expectedBody: "DELETETest",
	},
	{
		method:       http.MethodGet,
		path:         "/test/1",
		expectedBody: "GETTestParam1",
	},
	{
		method:       http.MethodPost,
		path:         "/test/1",
		expectedBody: "POSTTestParam1",
	},
	{
		method:       http.MethodPut,
		path:         "/test/1",
		expectedBody: "PUTTestParam1",
	},
	{
		method:       http.MethodDelete,
		path:         "/test/1",
		expectedBody: "DELETETestParam1",
	},
	{
		method:       http.MethodGet,
		path:         "/test/invalid/path",
		expectedBody: "NotFound",
	},
	{
		method:       http.MethodPost,
		path:         "/test/invalid/path",
		expectedBody: "NotFound",
	},
	{
		method:       http.MethodPut,
		path:         "/test/invalid/path",
		expectedBody: "NotFound",
	},
	{
		method:       http.MethodDelete,
		path:         "/test/invalid/path",
		expectedBody: "NotFound",
	},
	{
		method:          http.MethodGet,
		path:            "/onlyPostAllowed",
		expectedBody:    "Method Not Allowed\n",
		expectedHeaders: map[string]string{"Allow": http.MethodPost},
	},
	{
		method:          http.MethodGet,
		path:            "/onlyPutAllowed",
		expectedBody:    "Method Not Allowed\n",
		expectedHeaders: map[string]string{"Allow": http.MethodPut},
	},
	{
		method:          http.MethodGet,
		path:            "/onlyDeleteAllowed",
		expectedBody:    "Method Not Allowed\n",
		expectedHeaders: map[string]string{"Allow": http.MethodDelete},
	},
	{
		method:          http.MethodPost,
		path:            "/onlyGetAllowed",
		expectedBody:    "Method Not Allowed\n",
		expectedHeaders: map[string]string{"Allow": http.MethodGet},
	},
	{
		method:          http.MethodPost,
		path:            "/onlyPutAllowed",
		expectedBody:    "Method Not Allowed\n",
		expectedHeaders: map[string]string{"Allow": http.MethodPut},
	},
	{
		method:          http.MethodPost,
		path:            "/onlyDeleteAllowed",
		expectedBody:    "Method Not Allowed\n",
		expectedHeaders: map[string]string{"Allow": http.MethodDelete},
	},
	{
		method:          http.MethodPut,
		path:            "/onlyGetAllowed",
		expectedBody:    "Method Not Allowed\n",
		expectedHeaders: map[string]string{"Allow": http.MethodGet},
	},
	{
		method:          http.MethodPut,
		path:            "/onlyPostAllowed",
		expectedBody:    "Method Not Allowed\n",
		expectedHeaders: map[string]string{"Allow": http.MethodPost},
	},
	{
		method:          http.MethodPut,
		path:            "/onlyDeleteAllowed",
		expectedBody:    "Method Not Allowed\n",
		expectedHeaders: map[string]string{"Allow": http.MethodDelete},
	},
	{
		method:          http.MethodDelete,
		path:            "/onlyGetAllowed",
		expectedBody:    "Method Not Allowed\n",
		expectedHeaders: map[string]string{"Allow": http.MethodGet},
	},
	{
		method:          http.MethodDelete,
		path:            "/onlyPostAllowed",
		expectedBody:    "Method Not Allowed\n",
		expectedHeaders: map[string]string{"Allow": http.MethodPost},
	},
	{
		method:          http.MethodDelete,
		path:            "/onlyPutAllowed",
		expectedBody:    "Method Not Allowed\n",
		expectedHeaders: map[string]string{"Allow": http.MethodPut},
	},
	{
		method:          http.MethodGet,
		path:            "/postPutDelteAllowed",
		expectedBody:    "Method Not Allowed\n",
		expectedHeaders: map[string]string{"Allow": "POST, PUT, DELETE"},
	},
}

func TestRouter_ServeHTTP(t *testing.T) {
	rtr := &Router{
		NotFound: func(res http.ResponseWriter, req *http.Request) { fmt.Fprint(res, "NotFound") },
		routes:   make(map[string]route),
	}
	rtr.GET("test", func(res http.ResponseWriter, req *http.Request) { fmt.Fprint(res, "GETTest") })
	rtr.POST("test", func(res http.ResponseWriter, req *http.Request) { fmt.Fprint(res, "POSTTest") })
	rtr.PUT("test", func(res http.ResponseWriter, req *http.Request) { fmt.Fprint(res, "PUTTest") })
	rtr.DELETE("test", func(res http.ResponseWriter, req *http.Request) { fmt.Fprint(res, "DELETETest") })

	rtr.GET("test/:id", func(res http.ResponseWriter, req *http.Request) {
		param, _ := GetParam(req.Context(), "id")
		fmt.Fprint(res, fmt.Sprintf("GETTestParam%s", param))
	})
	rtr.POST("test/:id", func(res http.ResponseWriter, req *http.Request) {
		param, _ := GetParam(req.Context(), "id")
		fmt.Fprint(res, fmt.Sprintf("POSTTestParam%s", param))
	})
	rtr.PUT("test/:id", func(res http.ResponseWriter, req *http.Request) {
		param, _ := GetParam(req.Context(), "id")
		fmt.Fprint(res, fmt.Sprintf("PUTTestParam%s", param))
	})
	rtr.DELETE("test/:id", func(res http.ResponseWriter, req *http.Request) {
		param, _ := GetParam(req.Context(), "id")
		fmt.Fprint(res, fmt.Sprintf("DELETETestParam%s", param))
	})

	rtr.GET("onlyGetAllowed", func(res http.ResponseWriter, req *http.Request) { fmt.Fprint(res, "GETAllowed") })
	rtr.POST("onlyPostAllowed", func(res http.ResponseWriter, req *http.Request) { fmt.Fprint(res, "POSTAllowed") })
	rtr.PUT("onlyPutAllowed", func(res http.ResponseWriter, req *http.Request) { fmt.Fprint(res, "PUTAllowed") })
	rtr.DELETE("onlyDeleteAllowed", func(res http.ResponseWriter, req *http.Request) { fmt.Fprint(res, "DELETEAllowed") })

	rtr.POST("postPutDelteAllowed", func(res http.ResponseWriter, req *http.Request) { fmt.Fprint(res, "postPutDelteAllowed") })
	rtr.PUT("postPutDelteAllowed", func(res http.ResponseWriter, req *http.Request) { fmt.Fprint(res, "postPutDelteAllowed") })
	rtr.DELETE("postPutDelteAllowed", func(res http.ResponseWriter, req *http.Request) { fmt.Fprint(res, "postPutDelteAllowed") })

	for _, test := range serveTests {
		req, _ := http.NewRequest(test.method, test.path, nil)
		res := httptest.NewRecorder()
		rtr.ServeHTTP(res, req)
		bodyBytes, _ := ioutil.ReadAll(res.Body)
		if string(bodyBytes) != test.expectedBody {
			t.Errorf("Expected body to be: %s, but got: %s. Is the wrong handler called?", test.expectedBody, string(bodyBytes))
		}
		if test.expectedHeaders != nil {
			for k, v := range test.expectedHeaders {
				if res.Header().Get(k) != v {
					t.Errorf("Expected header: %s to be: %s, but was: %s", k, v, res.Header().Get(k))
				}
			}
		}
	}
}

var matchTests = []struct {
	// Route setup
	routePath string

	// Input
	inputPath string

	// Expected output
	match      bool
	parameters map[string]string
}{
	{
		routePath: "/",
		inputPath: "/",
		match:     true,
	},
	{
		routePath: "/test",
		inputPath: "/",
		match:     false,
	},
	{
		routePath: "/",
		inputPath: "/test",
		match:     false,
	},
	{
		routePath: "/test",
		inputPath: "/test/",
		match:     true,
	},
	{
		routePath: "/test/",
		inputPath: "/test2",
		match:     false,
	},
	{
		routePath: "/test2",
		inputPath: "/test",
		match:     false,
	},
	{
		routePath: "/test/test2/",
		inputPath: "/test/test2/",
		match:     true,
	},
	{
		routePath: "/test2/test",
		inputPath: "/test/test2",
		match:     false,
	},
	{
		routePath: "/test/test2",
		inputPath: "/test/",
		match:     false,
	},
	{
		routePath: "/test/test2/test3",
		inputPath: "/test/test2",
		match:     false,
	},
	{
		routePath: "/test/test2/test3",
		inputPath: "/test/",
		match:     false,
	},
	{
		routePath: "/test/test3/test3",
		inputPath: "/test/test2/test3",
		match:     false,
	},
	{
		routePath: "/group/:grpId/user/:usrId",
		inputPath: "/group/4/user/",
		match:     false,
	},
	{
		routePath: "/user/:usrId",
		inputPath: "/user/8",
		match:     true,
		parameters: map[string]string{
			"usrId": "8",
		},
	},
	{
		routePath: "/group/:grpId",
		inputPath: "/user/8",
		match:     false,
		parameters: map[string]string{
			"usrId": "8",
		},
	},
	{
		routePath: "/:id/user",
		inputPath: "/test/user",
		match:     true,
		parameters: map[string]string{
			"id": "test",
		},
	},
	{
		routePath: "/:id/user",
		inputPath: "/test/grp",
		match:     false,
		parameters: map[string]string{
			"id": "test",
		},
	},
	{
		routePath: "/group/:grpId/user/:usrId",
		inputPath: "/group/4/user/8",
		match:     true,
		parameters: map[string]string{
			"grpId": "4",
			"usrId": "8",
		},
	},
	{
		routePath: "/:id1/:id2/:id3/:id4",
		inputPath: "/test/testing/10/ten",
		match:     true,
		parameters: map[string]string{
			"id1": "test",
			"id2": "testing",
			"id3": "10",
			"id4": "ten",
		},
	},
}

func TestRouter_routeMatch(t *testing.T) {
	// Test for all combinations of http method on route and in request.
	for _, test := range matchTests {
		route := route{
			segments: segment(test.routePath),
		}

		ctx, match := route.match(context.Background(), test.inputPath)

		if test.match != match {
			t.Errorf("Expected match: %t, but got %t for test with routePath: %s", test.match, match, test.routePath)
		}

		if (test.parameters != nil) && (ctx != nil) {
			for k, v := range test.parameters {
				if val, exists := GetParam(ctx, k); !exists || val != v {
					t.Errorf("Expected param: %s to be %s, but was %s for test with routePath: %s", k, v, val, test.routePath)
				}
			}
		}

	}
}
