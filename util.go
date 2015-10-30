package mohttp

import (
	"fmt"
	"net/http"
)

func recoverErr() error {
	r := recover()

	if r == nil {
		return nil
	}

	if err, ok := r.(error); ok {
		return err
	}

	return fmt.Errorf("%#v", r)
}

func Redirect(path string) Handler {
	return HandlerFunc(func(c *Context) {
		http.Redirect(c.ResponseWriter(), c.Request(), path, http.StatusTemporaryRedirect)
		c.Next().Handle(c)
	})
}
