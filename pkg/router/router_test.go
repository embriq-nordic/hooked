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
	method       string
	path         string
	expectedBody string
	parameters   map[string]string
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
}

func TestRouter_ServeHTTP(t *testing.T) {
	rtr := &Router{
		NotFound: func(res http.ResponseWriter, req *http.Request) { fmt.Fprint(res, "NotFound") },
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

	for _, test := range serveTests {
		req, _ := http.NewRequest(test.method, test.path, nil)
		res := httptest.NewRecorder()
		rtr.ServeHTTP(res, req)
		bodyBytes, _ := ioutil.ReadAll(res.Body)
		if string(bodyBytes) != test.expectedBody {
			t.Errorf("Expected body to be: %s, but got: %s. Is the wrong handler called?", string(bodyBytes), test.expectedBody)
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

var supportedMethods = []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions}

func TestRouter_routeMatch(t *testing.T) {
	// Test for all combinations of http method on route and in request.
	for _, routeMethod := range supportedMethods {
		for _, inputMethod := range supportedMethods {
			for _, test := range matchTests {
				route := route{
					method:   routeMethod,
					segments: segment(test.routePath),
				}

				ctx, match := route.match(context.Background(), inputMethod, test.inputPath)

				// If method matches, evaluate the excpected result. If method doesnt match the match variable should always be false.
				if routeMethod == inputMethod {
					if test.match != match {
						t.Errorf("Expected match: %t, but got %t for test with routePath: %s and routeMethod: %s", test.match, match, test.routePath, routeMethod)
					}

					if (test.parameters != nil) && (ctx != nil) {
						for k, v := range test.parameters {
							if val, exists := GetParam(ctx, k); !exists || val != v {
								t.Errorf("Expected param: %s to be %s, but was %s for test with routePath: %s and routeMethod: %s", k, v, val, test.routePath, routeMethod)
							}
						}
					}
				} else {
					if match {
						t.Errorf("Methods disnt match. Expected match to be false but was true or test with routePath: %s and routeMethod: %s", test.routePath, routeMethod)
					}
				}
			}
		}
	}
}
