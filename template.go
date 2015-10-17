package http

import (
	"golang.org/x/net/context"
	"html/template"
)

const templateKey = contextKey("github.com/jonasi/http.Template")

type Template struct {
	Template *template.Template
}

func (t *Template) Handle(c *Context) {
	c.Context = context.WithValue(c.Context, templateKey, t)
	c.Next.Handle(c)
}

func TemplateResponse(c *Context, name string, data interface{}) {
	t := c.Context.Value(templateKey).(*Template)
	t.Template.ExecuteTemplate(c.Writer, name, data)
}
