package middleware

import (
	"encoding/json"
	"github.com/jonasi/mohttp"
	"golang.org/x/net/context"
)

var jsonContextValue = mohttp.NewContextValueStore("github.com/jonasi/http.JSON")

type JSONOptions struct {
	HandleErr func(context.Context, error) interface{}
	Transform func(interface{}) interface{}
}

func (j *JSONOptions) Handle(c context.Context) {
	c = jsonContextValue.Set(c, j)
	c = mohttp.WithResponder(c, &jsonResponder{j})

	mohttp.Next(c)
}

type jsonResponder struct {
	opts *JSONOptions
}

func (j *jsonResponder) HandleResult(c context.Context, data interface{}) error {
	if j.opts != nil && j.opts.Transform != nil {
		data = j.opts.Transform(data)
	}

	mohttp.GetResponseWriter(c).Header().Add("Content-Type", "application/json")
	return json.NewEncoder(mohttp.GetResponseWriter(c)).Encode(data)
}

func (j *jsonResponder) HandleErr(c context.Context, err error) {
	if h, ok := err.(*mohttp.HTTPError); ok {
		mohttp.GetResponseWriter(c).WriteHeader(h.Code)
	}

	if j.opts != nil && j.opts.HandleErr != nil {
		data := j.opts.HandleErr(c, err)
		err = j.HandleResult(c, data)
	}

	if err != nil {
		panic("No error json handler specified")
	}
}

func JSONHandler(fn mohttp.DataHandlerFunc) mohttp.Handler {
	return mohttp.HandlerFunc(func(c context.Context) {
		_, ok := jsonContextValue.Get(c).(*JSONOptions)

		if !ok {
			c = mohttp.WithResponder(c, &jsonResponder{})
		}

		fn(c)
	})
}

func JSONBodyDecode(c context.Context, dest interface{}) error {
	return json.NewDecoder(mohttp.GetRequest(c).Body).Decode(dest)
}
