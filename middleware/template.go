package middleware

import (
	"github.com/jonasi/mohttp"
	"golang.org/x/net/context"
	"html/template"
)

var templateContextValue = mohttp.NewContextValueStore("github.com/jonasi/http.templateHandler")

func Template(t *template.Template) mohttp.Handler {
	return &templateOptions{t}
}

type templateOptions struct {
	template *template.Template
}

func (t *templateOptions) Handle(c context.Context) {
	c = templateContextValue.Set(c, t)
	mohttp.Next(c)
}

func TemplateHandler(fn func(c context.Context) (string, map[string]interface{})) mohttp.Handler {
	return mohttp.HandlerFunc(func(c context.Context) {
		var (
			name, data = fn(c)
			t          = templateContextValue.Get(c).(*templateOptions)
		)
		t.template.ExecuteTemplate(mohttp.GetResponseWriter(c), name, data)
	})
}
