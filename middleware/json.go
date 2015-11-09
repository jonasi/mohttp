package middleware

import (
	"encoding/json"

	"github.com/jonasi/mohttp"
	"golang.org/x/net/context"
)

type JSONResponder struct {
	OnError   func(context.Context, error) interface{}
	Transform func(interface{}) interface{}
}

func (j *JSONResponder) HandleResult(c context.Context, data interface{}) error {
	if j.Transform != nil {
		data = j.Transform(data)
	}

	mohttp.GetResponseWriter(c).Header().Add("Content-Type", "application/json")
	return json.NewEncoder(mohttp.GetResponseWriter(c)).Encode(data)
}

func (j *JSONResponder) HandleErr(c context.Context, err error) {
	if h, ok := err.(mohttp.HTTPError); ok {
		mohttp.GetResponseWriter(c).WriteHeader(int(h))
	}

	if j.OnError != nil {
		data := j.OnError(c, err)
		err = j.HandleResult(c, data)
	}

	if err != nil {
		panic("No error json handler specified")
	}
}

func JSONHandler(fn mohttp.DataHandlerFunc) mohttp.Handler {
	return mohttp.HandlerFunc(func(c context.Context) {
		r, _ := mohttp.GetDataResponder(c)

		if _, ok := r.(*JSONResponder); !ok {
			c = mohttp.WithDataResponder(c, &JSONResponder{})
		}

		fn.Handle(c)
	})
}

func JSONBodyDecode(c context.Context, dest interface{}) error {
	r := mohttp.GetRequest(c)
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(dest)
}
