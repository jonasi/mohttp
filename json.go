package mohttp

import (
	"encoding/json"
)

var jsonContextValue = NewContextValueStore("github.com/jonasi/http.JSON")

var DefaultJSON = JSON(DefaultJSONTransform)

func DefaultJSONTransform(data interface{}) interface{} {
	if err, ok := data.(error); ok {
		return map[string]interface{}{
			"error": err,
		}
	}

	return data
}

func JSON(fn func(interface{}) interface{}) Handler {
	return &jsonHandler{fn}
}

type jsonHandler struct {
	TransformFunc func(interface{}) interface{}
}

func (j *jsonHandler) Handle(c *Context) {
	c = jsonContextValue.Set(c, j)
	c.Writer.Header().Add("Content-Type", "application/json")

	defer func() {
		if err := recoverErr(); err != nil {
			j.handleErr(c, err)
		}
	}()

	c.Next.Handle(c)
}

func (j *jsonHandler) handleErr(c *Context, err error) {
	b, _ := json.Marshal(map[string]string{
		"error": err.Error(),
	})

	c.Writer.Write(b)
}

func JSONResponse(c *Context, data interface{}) {
	j := jsonContextValue.Get(c).(*jsonHandler)

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
