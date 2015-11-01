package mohttp

import (
	"encoding/json"
	"golang.org/x/net/context"
)

var jsonContextValue = NewContextValueStore("github.com/jonasi/http.JSON")

type JSONOptions struct {
	HandleErr func(context.Context, error) interface{}
	Transform func(interface{}) interface{}
}

func (j *JSONOptions) Handle(c context.Context) {
	c = jsonContextValue.Set(c, j)
	c = WithResponder(c, &jsonResponder{j})

	GetNext(c).Handle(c)
}

type jsonResponder struct {
	opts *JSONOptions
}

func (j *jsonResponder) HandleResult(c context.Context, data interface{}) error {
	if j.opts != nil && j.opts.Transform != nil {
		data = j.opts.Transform(data)
	}

	GetResponseWriter(c).Header().Add("Content-Type", "application/json")
	return json.NewEncoder(GetResponseWriter(c)).Encode(data)
}

func (j *jsonResponder) HandleErr(c context.Context, err error) {
	if j.opts != nil && j.opts.HandleErr != nil {
		data := j.opts.HandleErr(c, err)
		err = j.HandleResult(c, data)
	}

	if err != nil {
		panic("No error json handler specified")
	}
}

func JSONHandler(fn DataHandlerFunc) Handler {
	h := DataHandler(fn)

	return HandlerFunc(func(c context.Context) {
		_, ok := jsonContextValue.Get(c).(*JSONOptions)

		if !ok {
			c = WithResponder(c, &jsonResponder{})
		}

		h.Handle(c)
	})
}
