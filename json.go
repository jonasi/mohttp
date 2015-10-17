package http

import (
	"encoding/json"
	"golang.org/x/net/context"
)

const jsonKey = contextKey("github.com/jonasi/http.JSON")

var DefaultJSON = &JSON{
	TransformFunc: DefaultJSONTransform,
}

func DefaultJSONTransform(data interface{}) interface{} {
	if err, ok := data.(error); ok {
		return map[string]interface{}{
			"error": err,
		}
	}

	return data
}

type JSON struct {
	TransformFunc func(interface{}) interface{}
}

func (j *JSON) Handle(c *Context) {
	c.Context = context.WithValue(c.Context, jsonKey, j)
	c.Writer.Header().Add("Content-Type", "application/json")

	defer func() {
		if err := recoverErr(); err != nil {
			j.handleErr(c, err)
		}
	}()

	c.Next.Handle(c)
}

func JSONResponse(c *Context, data interface{}) {
	j := c.Context.Value(jsonKey).(*JSON)

	if j.TransformFunc != nil {
		data = j.TransformFunc(data)
	}

	b, err := json.Marshal(data)

	if err != nil {
		j.handleErr(c, err)
		return
	}

	c.Writer.Write(b)
}

func (j *JSON) handleErr(c *Context, err error) {
	b, _ := json.Marshal(map[string]string{
		"error": err.Error(),
	})

	c.Writer.Write(b)
}
