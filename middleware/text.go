package middleware

import (
	"golang.org/x/net/context"
)

type TextHandlerFunc func(context.Context) (string, error)

func (fn TextHandlerFunc) Handle(c context.Context) {

}
