package lambdahandler

import (
	"bytes"
	"context"
	"github.com/aws/aws-lambda-go/events"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Handler is the entry point for a lambda and wraps the APIGateway events so http.Handler funcs can be called.
type Handler struct {
	Handler http.Handler
}

// Handle wraps the APIGatewayProxyRequest/Response so we can use regular http handler funcs.
func (h Handler) Handle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	httpRes := lambdaResponseWriter{
		headers: make(http.Header),
		buffer:  &bytes.Buffer{},
		status:  http.StatusOK,
	}

	httpReq := &http.Request{
		Method:        req.HTTPMethod,
		URL:           &url.URL{Path: req.Path},
		Header:        req.MultiValueHeaders,
		Body:          ioutil.NopCloser(bytes.NewBufferString(req.Body)),
		ContentLength: int64(len(req.Body)),
		RemoteAddr:    req.RequestContext.Identity.SourceIP,
	}

	h.Handler.ServeHTTP(&httpRes, httpReq)

	payload, err := ioutil.ReadAll(httpRes.buffer)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode:        httpRes.status,
		MultiValueHeaders: httpRes.headers,
		Body:              string(payload),
	}, nil
}

type lambdaResponseWriter struct {
	headers http.Header
	buffer  io.ReadWriter
	status  int
}

func (l *lambdaResponseWriter) Header() http.Header {
	return l.headers
}

func (l *lambdaResponseWriter) Write(p []byte) (int, error) {
	return l.buffer.Write(p)
}

func (l *lambdaResponseWriter) WriteHeader(statusCode int) {
	l.status = statusCode
}
