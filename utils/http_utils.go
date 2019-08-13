package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"testing"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
	"gotest.tools/assert"
)

func WriteJSONOrError(obj interface{}, ctx *fasthttp.RequestCtx, err error, msg string, statusCode int) {
	if err != nil {
		ctx.Error(msg, statusCode)
	} else {
		WriteJSON(obj, ctx)
	}
}

func WriteJSON(obj interface{}, ctx *fasthttp.RequestCtx) {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		ctx.Error("Internal Error: Marshaling", http.StatusInternalServerError)
	} else {
		ctx.Response.Header.SetContentType("application/json;cahrset=utf8")
		ctx.Write(jsonBytes)
	}
}

func ApiClient(netListener *fasthttputil.InmemoryListener, router *fasthttprouter.Router) (*http.Client, error) {
	var err error
	go func() {
		err = fasthttp.Serve(netListener, router.Handler)
	}()

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return netListener.Dial()
			},
		},
	}

	return client, err
}

func DO_POST(url string, obj interface{}, expectedStatus int, client *http.Client, t *testing.T) *http.Response {
	return DO_HTTP_Req("POST", url, obj, expectedStatus, client, t)
}

func DO_GET(url string, expectedStatus int, client *http.Client, t *testing.T) *http.Response {
	return DO_HTTP_Req("GET", url, nil, expectedStatus, client, t)
}

func DO_HTTP_Req(method string, url string, obj interface{}, expectedStatus int, client *http.Client, t *testing.T) *http.Response {

	var buffer io.Reader
	if obj != nil {
		jsonBytes, err := json.Marshal(obj)
		assert.NilError(t, err)
		buffer = bytes.NewBuffer(jsonBytes)
	}
	req, err := http.NewRequest(method, url, buffer)
	assert.NilError(t, err, "Failed while creating request")
	resp, err := client.Do(req)
	assert.NilError(t, err, "Failed to execute request")
	assert.Equal(t, expectedStatus, resp.StatusCode)
	assert.Equal(t, "application/json;cahrset=utf8", resp.Header.Get("content-type"))

	return resp
}
