package mohttp

import (
	"golang.org/x/net/context"
	"html/template"
)

var templateContextValue = NewContextValueStore("github.com/jonasi/http.templateHandler")

func Template(t *template.Template) Handler {
	return &templateOptions{t}
}

type templateOptions struct {
	template *template.Template
}

func (t *templateOptions) Handle(c context.Context) {
	c = templateContextValue.Set(c, t)
	Next(c)
}

func TemplateHandler(fn func(c context.Context) (string, map[string]interface{})) Handler {
	return HandlerFunc(func(c context.Context) {
		var (
			name, data = fn(c)
			t          = templateContextValue.Get(c).(*templateOptions)
		)
		t.template.ExecuteTemplate(GetResponseWriter(c), name, data)
	})
}
