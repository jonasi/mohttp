package mohttp

import (
	"html/template"
)

var templateContextValue = NewContextValueStore("github.com/jonasi/http.templateHandler")

func Template(t *template.Template) Handler {
	return &templateHandler{t}
}

type templateHandler struct {
	template *template.Template
}

func (t *templateHandler) Handle(c *Context) {
	c = templateContextValue.Set(c, t)
	c.Next().Handle(c)
}

func TemplateResponse(c *Context, name string, data interface{}) {
	t := templateContextValue.Get(c).(*templateHandler)
	t.template.ExecuteTemplate(c.ResponseWriter(), name, data)
}
