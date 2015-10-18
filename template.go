package http

import (
	"golang.org/x/net/context"
	"html/template"
)

const templateKey = contextKey("github.com/jonasi/http.templateHandler")

func Template(t *template.Template) Handler {
	return &templateHandler{t}
}

type templateHandler struct {
	template *template.Template
}

func (t *templateHandler) Handle(c *Context) {
	c.Context = context.WithValue(c.Context, templateKey, t)
	c.Next.Handle(c)
}

func TemplateResponse(c *Context, name string, data interface{}) {
	t := c.Context.Value(templateKey).(*templateHandler)
	t.template.ExecuteTemplate(c.Writer, name, data)
}
